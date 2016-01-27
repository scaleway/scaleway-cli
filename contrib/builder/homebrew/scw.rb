require "language/go"

class Scw < Formula
  desc "Manage BareMetal Servers from Command Line (as easily as with Docker)"
  homepage "https://github.com/scaleway/scaleway-cli"
  url "https://github.com/scaleway/scaleway-cli/archive/v1.7.0.tar.gz"
  sha256 "c1afa976301e8b2e8ce6b4c4f306cbfee5171f3a17f3083149e338c2ada774a5"

  head "https://github.com/scaleway/scaleway-cli.git"

  depends_on "go" => :build

  go_resource "github.com/Sirupsen/logrus" do
    url "https://github.com/Sirupsen/logrus.git", revision: "f7f79f729e0fbe2fcc061db48a9ba0263f588252"
  end

  def install
    ENV["GOPATH"] = buildpath
    ENV["CGO_ENABLED"] = "0"
    ENV["GO15VENDOREXPERIMENT"] = "1"
    ENV.prepend_create_path "PATH", buildpath/"bin"

    mkdir_p buildpath/"src/github.com/scaleway"
    ln_s buildpath, buildpath/"src/github.com/scaleway/scaleway-cli"
    Language::Go.stage_deps resources, buildpath/"src"

    version = `cat .goxc.json | grep "PackageVersion" | egrep -o "([0-9]{1,}\.)+[0-9]{1,}"`
    system "go", "build", "-ldflags", "-X  github.com/scaleway/scaleway-cli/pkg/scwversion.GITCOMMIT=homebrew -X  github.com/scaleway/scaleway-cli/pkg/scwversion.VERSION=#{version}", "./cmd/scw"
    bin.install "scw"

    bash_completion.install "contrib/completion/bash/scw"
    zsh_completion.install "contrib/completion/zsh/_scw"
  end

  def caveats
    "Use `scw login` to set up the correct environment to use scaleway-cli."
  end

  test do
    output = shell_output(bin/"scw version")
    assert output.include? "OS/Arch (client): darwin/amd64"
  end
end
