# Frequently Asked Questions

## General

**What is Spark?**
Spark is a global creator platform for live streaming, video content, and digital interaction. It connects creators with audiences worldwide.

**Is Spark open source?**
Core platform components are open source under the MIT license. Certain proprietary features (e.g., advanced ML models, enterprise compliance) are closed source.

**Who owns the content I upload?**
You retain full ownership of your content. Spark is granted a license to host, transcode, and deliver your content on the platform.

## Development

**How do I set up the development environment?**
See `README.md` for quick start instructions and `docs/getting-started.md` for a detailed walkthrough.

**What package manager do you use?**
We use **pnpm** with workspaces for the monorepo.

**How do I add a new service?**
Create a new directory under `services/`, add a `package.json` and Dockerfile, register it in the root `pnpm-workspace.yaml`, and submit a PR.

**How do I run tests?**
Use `pnpm test` for all packages or `pnpm --filter <package> test` for a specific package.

## Platform

**Can I stream from a mobile device?**
Yes. The Flutter mobile app supports live streaming from iOS and Android devices.

**What video formats are supported?**
We accept MP4, MOV, AVI, and WebM for uploads. Live streaming uses WebRTC (browser) or RTMP (encoder). Output is HLS with adaptive bitrate.

**How do I report a bug?**
Open an issue on GitHub with the bug report template, or email support@sparkplatform.com.
