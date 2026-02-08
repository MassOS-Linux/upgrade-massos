# upgrade-massos
Modern system upgrade utility for MassOS

# About
`upgrade-massos` is designed to be a replacement for the older `massos-upgrade`
utility (which is now discontinued). Unlike `massos-upgrade`, `upgrade-massos`
provides the following features which were missing from the original:

- Support for MassOS experimental and 2026+ builds.
- Support for an improved online update server structure (described below).
- Support for checksum and GPG signature validation for online updates.
- Support for upgrading to a build with a newer major Glibc version.
- Support for cleaning up stale files from the older version after the upgrade.

As well as, of course, retaining the following core features from the original:

- Support for upgrading while the system is running.
- Support for preserving user-customized configuration files under `/etc`.

# Compatibility
MassOS **experimental-20251105** or higher. If you are running an older build,
please back up your personal data and perform a clean reinstallation of the
latest available build. Then, subsequent updates can be performed with this
utility.

**NOTE:** Upgrading from any pre-2023 version, including **2022.10** or older,
is entirely **UNSUPPORTED**.

# Running
The utility is expected to be pre-installed in MassOS builds from March 2026
and onwards. If so, you can simply invoke it by running one of the following
commands:
```sh
# To check online for a new update and install it if available:
sudo upgrade-massos
```

```sh
# To install an update offline, from a local rootfs image:
sudo upgrade-massos <rootfs-file>.tar.zst
```

If your system does **NOT** have `upgrade-massos` pre-installed, due to being
an older build which is still supported by the utility, then you can download
and run it from the git source tree with the following commands:
```sh
# Download the utility and change to its directory:
git clone https://github.com/MassOS-Linux/upgrade-massos
cd upgrade-massos
```

And then run the utility as follows:
```sh
# To check online for a new update and install it if available:
sudo ./upgrade-massos
```

```sh
# To install an update offline, from a local rootfs image:
sudo ./upgrade-massos <rootfs-file>.tar.zst
```

If a new update was found (online), or your local rootfs image is valid
(offline), then you will be prompted if you want to start the upgrade. Answer
`y` to begin the upgrade process.

Please note that the upgrade process may take some time. Do **NOT** attempt to
run any other programs on your system while the upgrade is in progress, and do
**NOT** shut down, restart or suspend the system until the upgrade is fully
complete.

Once the upgrade process completes, you will be prompted to restart your
system. It is strongly recommended that you answer `y`, since it is important
to reboot as soon as possible after an upgrade. If you answer `n`, then please
manually restart your system as soon as is convenient for you to do so.

# Configuration
You can configure the utility by editing the file
`/etc/upgrade-massos/upgrade-massos.conf` (or simply `upgrade-massos.conf` if
you are running directly from the git source tree checkout).

TBC

# Server
TBC
