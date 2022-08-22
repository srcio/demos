# goreleaser
Go 程序项目的自动化发布工具，简单的发布命令帮助我们省去大量的重复工作。
## 安装
**MacOs**
```bash
brew install goreleaser/tap/goreleaser
```
**源码编译**
```bash
git clone https://github.com/goreleaser/goreleaser
cd goreleaser
go get ./...
go build -o goreleaser .
./goreleaser --version
```
## 初始化
在go项目下运行以下命令，生成`.goreleaser.yml`文件：
```bash
goreleaser init
```
生成文件后，自行配置，相信你能看的懂。
### 验证 .goreleaser.yml
```bash
goreleaser check
```
使用本地环境构建
```bash
goreleaser build --single-target
```
## 配置 github token
`token` 必须至少包含 `write:package` 权限，才能上传到发布资源中。

从github生成token，写入文件：
```bash
mkdir ~/.config/goreleaser
vim ~/.config/goreleaser/github_token
```
或者直接在终端导入环境配置：
```bash
export GITHUB_TOKEN="YOUR_GITHUB_TOKEN"
```
## 为项目打上 tag
```bash
git tag v0.1.0 -m "release v0.1.0"
git push origin v0.1.0
```

## 执行 goreleaser

```bash
goreleaser --rm-dist

# 跳过 git 修改验证
goreleaser build --skip-validate --rm-dist
```
如果目前还不想打tag，可以基于最近的一次提交 直接使用一下 `gorelease` 命令：
```bash
goreleaser build --snapshot
```
## Dry run
如果你想在真实的发布之前做发布测试，你可以尝试 `dry run` 。
### Build only 模式
使用如下命令，只会编译当前项目，可以验证项目是否存储编译错误。
```bash
goreleaser build
```
### Release 标识
使用 `--skip-publish` 命令标识来跳过发布到远端。
```bash
goreleaser release --skip-publish
```
```bash
goreleaser release --snapshot --skip-publish --rm-dist
```
## 工作原理
条件：
- 项目下有一个规范的 `·goreleaser.yml` 文件
- 干净的 `git` 目录
- 符合 `SemVer-compatible` 的版本命名
步骤：
1. defaulting：为每个步骤配置合理的默认值
2. building：构建二进制文件（binaries），归档文件（archives），包文件（packages），Docker镜像等
3. publishing：将构建的文件发布到配置的 GitHub Release 资源，Docker Hub 等
4. annoucing：发布完成后，配置通知

步骤可以通过命令标识 `--skip-{step_name}` 来跳过，如果当中某步骤失败，之后的步骤将不会运行。
## CGO
很遗憾，如果你的 Go 程序项目需要使用到 CGO 的跨平台编译，Docker 镜像是不支持的，并且你的配置将不会看到 `clean` 。

也许这里可以帮助到你：[Cross-compiling Go with CGO - GoReleaser](https://goreleaser.com/cookbooks/cgo-and-crosscompiling/)