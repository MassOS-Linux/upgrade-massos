/*
	ugm-install-helper - A file installation helper for upgrade-massos.

	Copyright (C) 2026 Daniel Massey. Covered under the same license as
	upgrade-massos itself (GPLv3 or later).

	This helper program exists to function as a wrapper for the bulk
	installation of files performed as part of the MassOS system upgrade
	process, and intends to significantly increase the speed of certain
	parts of the upgrade process.

	The fastest ways to recursively copy a tree of files to a destination,
	while preserving ownership and permissions, are 'cp -a' and 'rsync'.
	Unfortunately, neither of these programs safely handle the process of
	overwriting shared libraries. As a result, if either of these programs
	replace a shared library that is currently in use by any running
	program (which is inevitable when you are upgrading all the libraries
	of a running system), the running program will crash. In the case of
	upgrade-massos, this may very likely crash the entire desktop, as well
	as terminal currently running the upgrade utility.

	The 'install' utility, by contrast, does take the necessary precautions
	to ensure shared libraries are overwritten safely. However, 'install'
	has no bulk-install functionality like 'cp', nor does it have the
	ability to preserve permissions and ownership. As a result, using it
	requires performing shell iterations, 'stat' invocations, as well as
	calling the 'install' program tens of thousands of times. Reloading a
	program back into memory introduces significant slowdowns. By contrast,
	iterating in a language that compiles to native machine code (Go in
	this case, but C/C++/Rust would also be applicable), and calling only
	its standard library methods, does NOT introduce slowdowns, as the
	standard library itself remains in memory for the entire duration of
	the program's duration.

	TL;DR: This program mimics the functionality and speed of 'cp -a', with
	the safety and robustness of 'install.

	This program also implements the progress percentage indicator that was
	previously in use by the main upgrade-massos utility during the stage
	took the longest. This now means that any stage making use of this
	helper will also now have a progress indicator.

	This program also installs files in parallel, up to 4 parallel jobs at
	once, depending on available CPU threads, to speed up the process,
	especially when installing hundreds or thousands of very small files.
	It will not go above 4 parallel jobs no matter how many threads the CPU
	has, otherwise disk thrashing could occur on HDDs and weaker SSDs.

	This program currently behaves only as necessary for upgrade-massos. It
	does NOT accept any command-line parameters other than <source> and
	<destination>, and it does NOT preserve XATTRs (as no package in the
	base MassOS system makes use of them). Eventually this code could
	therefore be expanded upon and put into a standalone program. Anyone
	who has the time and desire to do that can "Go" right ahead (get it?)!
*/

package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"sync/atomic"
	"syscall"
)

type job struct {
	src string
	dst string
}

var totalFiles int64
var copiedFiles int64

func copyFileSafely(src, dst string, mode os.FileMode, uid, gid int) error {
	dir := filepath.Dir(dst)

	tmp, err := os.CreateTemp(dir, ".tmp-install-*")
	if err != nil {
		return err
	}
	tmpName := tmp.Name()

	defer func() {
		tmp.Close()
		os.Remove(tmpName)
	}()

	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	if _, err := io.Copy(tmp, in); err != nil {
		return err
	}

	if err := tmp.Chown(uid, gid); err != nil {
		return err
	}
	if err := tmp.Chmod(mode); err != nil {
		return err
	}
	if err := tmp.Sync(); err != nil {
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}

	return os.Rename(tmpName, dst)
}

func printProgress() {
	done := atomic.LoadInt64(&copiedFiles)
	percent := (float64(done) / float64(totalFiles)) * 100
	fmt.Printf("\rProgress: %.1f%%", percent)
}

func worker(jobs <-chan job, wg *sync.WaitGroup) {
	defer wg.Done()

	for j := range jobs {
		info, err := os.Lstat(j.src)
		if err != nil {
			fmt.Fprintln(os.Stderr, "\nstat failed:", err)
			continue
		}

		stat := info.Sys().(*syscall.Stat_t)
		uid := int(stat.Uid)
		gid := int(stat.Gid)

		switch {
		case info.Mode()&os.ModeSymlink != 0:
			target, err := os.Readlink(j.src)
			if err == nil {
				os.Remove(j.dst)
				if err := os.Symlink(target, j.dst); err == nil {
					_ = os.Lchown(j.dst, uid, gid)
				}
			}

		case info.Mode().IsRegular():
			// Take care to preserve special permission bits.
			mode := info.Mode() & (os.ModePerm | os.ModeSetuid | os.ModeSetgid | os.ModeSticky)
			_ = copyFileSafely(j.src, j.dst, mode, uid, gid)
		}

		atomic.AddInt64(&copiedFiles, 1)
		printProgress()
	}
}

func main() {
	if len(os.Args) != 3 {
		fmt.Fprintf(os.Stderr, "usage: %s <source> <destination>\n", os.Args[0])
		os.Exit(1)
	}

	srcRoot := filepath.Clean(os.Args[1])
	dstParent := filepath.Clean(os.Args[2])

	base := filepath.Base(srcRoot)
	dstRoot := filepath.Join(dstParent, base)

	// Create all directories first, for robustness.
	var jobsList []job

	err := filepath.WalkDir(srcRoot, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		rel, _ := filepath.Rel(srcRoot, path)
		dstPath := filepath.Join(dstRoot, rel)

		info, err := d.Info()
		if err != nil {
			return err
		}

		stat := info.Sys().(*syscall.Stat_t)
		uid := int(stat.Uid)
		gid := int(stat.Gid)

		if info.IsDir() {
			// Take care to preserve special permission bits.
			_ = os.MkdirAll(dstPath, 0755)
			if err := os.Chown(dstPath, uid, gid); err != nil {
				return err
			}
			// Only chmod now we now special bits won't be lost.
			mode := info.Mode() & (os.ModePerm | os.ModeSetuid | os.ModeSetgid | os.ModeSticky)
			if err := os.Chmod(dstPath, mode); err != nil {
				return err
			}
			return nil
		}

		if info.Mode().IsRegular() || (info.Mode()&os.ModeSymlink != 0) {
			jobsList = append(jobsList, job{src: path, dst: dstPath})
		}
		return nil
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, "walk failed:", err)
		os.Exit(1)
	}

	totalFiles = int64(len(jobsList))
	if totalFiles == 0 {
		fmt.Println("Nothing to install.")
		return
	}

	fmt.Printf("Progress: 0.0%%")

	// Start parallel workers (maximum of 4 jobs).
	numWorkers := min(runtime.NumCPU(), 4)
	jobs := make(chan job, numWorkers)
	var wg sync.WaitGroup

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go worker(jobs, &wg)
	}

	for _, j := range jobsList {
		jobs <- j
	}
	close(jobs)

	wg.Wait()
	fmt.Printf("\rProgress: 100.0%%\n")
}
