require "language/go"

class Scw < Formula
  homepage "https://github.com/scaleway/scaleway-cli"
  url "https://github.com/scaleway/scaleway-cli/archive/v1.1.0.tar.gz"
  sha256 "571146c067da9d3f411fabd24adb08f409f913e2237d80233a71655598b32a37"

  head "https://github.com/scaleway/scaleway-cli.git"

  depends_on "go" => :build
  depends_on "hg" => :build

  go_resource "github.com/tools/godep" do
    url "https://github.com/tools/godep.git", :revision => "58d90f262c13357d3203e67a33c6f7a9382f9223"
  end

  go_resource "github.com/kr/fs" do
    url "https://github.com/kr/fs.git", :revision => "2788f0dbd16903de03cb8186e5c7d97b69ad387b"
  end

  go_resource "golang.org/x/tools" do
    url "https://github.com/golang/tools.git", :revision => "473fd854f8276c0b22f17fb458aa8f1a0e2cf5f5"
  end

  go_resource "golang.org/x/crypto" do
    url "https://github.com/golang/crypto.git", :revision => "8b27f58b78dbd60e9a26b60b0d908ea642974b6d"
  end

  go_resource "github.com/scaleway/scaleway-cli" do
    url "https://github.com/scaleway/scaleway-cli.git", :revision => "49587ef360980cd28e7e0eac30806fb66372f2ed"
  end


  def install
    ENV["GOPATH"] = buildpath
    ENV["CGO_ENABLED"] = "0"
    ENV.prepend_create_path "PATH", buildpath/"bin"

    mkdir_p buildpath/"src/github.com/scaleway"
    ln_s buildpath, buildpath/"src/github.com/scaleway/scaleway-cli"
    Language::Go.stage_deps resources, buildpath/"src"

    cd "src/github.com/tools/godep" do
      system "go", "install"
    end

    system "make scwversion/version.go"
    system "./bin/godep", "get"
    system "./bin/godep", "go", "build", "-o", "scw"
    bin.install "scw"

    bash_completion.install "contrib/completion/bash/scw"
    zsh_completion.install "contrib/completion/zsh/_scw"
  end

  test do
    output = shell_output(bin/"scw version")
    assert output.include? "Client version: v1.1.0"
    assert output.include? "Git commit (client): 49587ef360980cd28e7e0eac30806fb66372f2ed"
  end
end
