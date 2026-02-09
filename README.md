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
you are running directly from the git source tree checkout). Comments in the
file describe the purpose of each option.

The upgrade channel can be specified via the `channel:` option. The available
channels will vary depending on the server, but the official MassOS upgrade
server will offer the following channels:

- `experimental` - Experimental builds of MassOS. During MassOS's experimental
  phase, this is the only channel currently available. After MassOS returns to
  regular releases, this channel will contain preview and testing builds.
- `stable` - Stable release builds of MassOS. This channel is not currently
  supported. After MassOS returns to normal releases and exits its experimental
  phase, this channel will be used for standard release upgrades, and will
  therefore be the recommended channel for most users. Stable releases will be
  well-tested and production ready, while unstable builds may encounter bugs
  or stability issues.

# Server
This section will describe the server structure expected by `upgrade-massos`.
It is only useful for individuals who wish to host their own upgrade server
(for example, in order to serve custom/unofficial MassOS builds).

This example layout is going to assume a server hosted at
`https://example.org/massos-builds`. The subdirectory used to store builds can
be anything you wish, in this example it is `massos-builds`.
```
massos-builds/
├── .upgrade-massos.yml
├── .upgrade-massos.CHANNELNAME1.yml
├── .upgrade-massos.CHANNELNAME2.yml
├── VERSIONCODE1
│   ├── massos-VERSIONCODE1-rootfs-x86_64-VARIANT.tar.zst
│   ├── massos-VERSIONCODE1-rootfs-x86_64-VARIANT.tar.zst.asc
│   └── massos-VERSIONCODE1-rootfs-x86_64-VARIANT.tar.zst.b2
└── VERSIONCODE2
    ├── massos-VERSIONCODE2-rootfs-x86_64-VARIANT.tar.zst
    ├── massos-VERSIONCODE2-rootfs-x86_64-VARIANT.tar.zst.asc
    └── massos-VERSIONCODE2-rootfs-x86_64-VARIANT.tar.zst.b2
```

The file `.upgrade-massos.yml` is the main server configuration information. It
identifies the server as a valid MassOS upgrade server, and dictates the
server's checksum/signature policies. It is structured as follows:
```
format: 1
checksum-policy: POLICY
signature-policy: POLICY
```
For `checksum-policy` and `signature-policy`, the following options in place of
`POLICY` are supported:

- `strict` - All builds must contain a valid checksum/signature, which must be
  validated by the client when performing an online update. If any build is
  missing the checksum/signature, or if the validation fails for any reason,
  the upgrade will NOT proceed.
- `sloppy` - If a build contains a valid checksum/signature, it/they will be
  validated by the client when performing an online update. If the validation
  fails for any reason, the upgrade will NOT proceed. If the checksum/signature
  file does NOT exist on the server, the upgrade WILL proceed anyway (skipping
  validation). This option is not recommended other than for testing purposes.
- `none` - No checksum/signature validation will be performed at all. This
  option is not recommended under any circumstance, and is subject to removal
  or deprecation in future updates to `upgrade-massos`.


Files named `.upgrade-massos.CHANNELNAME.yml` must exist for every channel that
the server supports. For example, `.upgrade-massos.experimental.yml` and
`.upgrade-massos.stable.yml`. You can choose whichever channel name(s) you wish
to use, but we recommend clearly documenting to users of your upgrade server
which channels are supported. The client must set the channel they wish to use
in their `upgrade-massos.conf` file. The channel file is structured as follows:
```
format: 1
version: LATESTVERSION
variants: [VARIANT1, VARIANT2]
compression: COMPRESSIONALGORITHM
checksum: CHECKSUMALGORITHM
signature: SIGNATUREFORMAT
```

Set (and update) `version` to the latest available version in this channel. It
is the version that the client will download and install when performing an
online update.

`variants` must contain an array of MassOS build variants supported by the
version. If there is only one supported, enter (e.g.) `[xfce]`. If multiple are
supported, write them as (e.g.) `[xfce, kde, gnome]`.

`compression` indicates the compression algorithm used by the build. It should
match the file extension after `.tar.`, e.g., `zst` (modern default for MassOS)
or `xz` (older builds or niche situations).

`checksum` indicates the file extension of the checksum file (and hence, the
checksum algorithm used). For example, `b2` indicates a Blake-2 checksum is
used, `sha256` indicates a SHA-256 checksum is used, and so on.

`signature` indicates the file extension of the build's detached GPG signature.
In most cases it will be `asc`, assuming the signature was produced using a
command such as `gpg --detach-sign --armor`.

Finally, the builds for each version should exist in subdirectories of the
version's code. For example, `experimental-20260208`.

At your own discretion, you can choose whether or not to (a) delete older
builds from your server after newer ones are made available, and (b) provide
ISO images alongside the rootfs tarballs. Both of these are unnecessary for
`upgrade-massos`, but may be useful if your server has a web interface for its
filesystem, and if your users may want to directly download builds from the
same server, rather than downloading them from a different source.

If you want to see an example server layout for reference, you may visit the
[official MassOS upgrade server](https://archive.massos.org/massos-builds/).
The official MassOS upgrade server's configuration files can be found
[here](https://archive.massos.org/massos-builds/.upgrade-massos.yml), and
[here](https://archive.massos.org/massos-builds/.upgrade-massos.experimental.yml).
