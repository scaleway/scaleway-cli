require "language/go"

class Scw < Formula
  homepage "https://github.com/scaleway/scaleway-cli"
  url "https://github.com/scaleway/scaleway-cli/archive/v1.0.0.tar.gz"
  sha256 "7a674320bc3298bf50dd2b50a9312c5baa159d5c2bf6f9abae539c9fc931fdce"

  head "https://github.com/scaleway/scaleway-cli.git"

  depends_on "go" => :build

  go_resource "github.com/tools/godep" do
    url "https://github.com/tools/godep.git", :revision => "58d90f262c13357d3203e67a33c6f7a9382f9223"
  end

  go_resource "github.com/kr/fs" do
    url "https://github.com/kr/fs.git", :revision => "2788f0dbd16903de03cb8186e5c7d97b69ad387b"
  end

  go_resource "golang.org/x/tools" do
    url "https://github.com/golang/tools.git", :revision => "473fd854f8276c0b22f17fb458aa8f1a0e2cf5f5"
  end

  #go_resource "golang.org/x/crypto" do
  #  url "https://github.com/golang/crypto.git", :revision => "8b27f58b78dbd60e9a26b60b0d908ea642974b6d"
  #end

  #go_resource "github.com/scaleway/scaleway-cli" do
  #  url "https://github.com/scaleway/scaleway-cli.git", :revision => "799e3bcdd7c72fdbb5d6dcf65f99654fc0fffc0b"
  #end


  def install
    ENV["GOPATH"] = buildpath
    ENV["CGO_ENABLED"] = "0"
    ENV.prepend_create_path "PATH", buildpath/"bin"

    #mkdir_p buildpath/"Godeps/_workspace/src/github.com/scaleway"
    #ln_s buildpath, buildpath/"Godeps/_workspace/src/github.com/scaleway/scaleway-cli"

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
    assert output.include? "Client version: v1.0.0"
    assert output.include? "Git commit (client): 799e3bcdd7c72fdbb5d6dcf65f99654fc0fffc0b"
  end
end
