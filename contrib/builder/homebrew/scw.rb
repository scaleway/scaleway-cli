require "language/go"

class Scw < Formula
  desc "Manage BareMetal Servers from Command Line (as easily as with Docker)"
  homepage "https://github.com/scaleway/scaleway-cli"
  url "https://github.com/scaleway/scaleway-cli/archive/v1.4.0.tar.gz"
  sha256 "ab7ee002be9557eb2b8075e3b0df340f5e379545152049f2512f1dc2b47b7b8a"

  head "https://github.com/scaleway/scaleway-cli.git"

  depends_on "go" => :build

  def install
    ENV["GOPATH"] = buildpath
    ENV["CGO_ENABLED"] = "0"
    ENV.prepend_create_path "PATH", buildpath/"bin"

    mkdir_p buildpath/"src/github.com/scaleway"
    ln_s buildpath, buildpath/"src/github.com/scaleway/scaleway-cli"
    Language::Go.stage_deps resources, buildpath/"src"

    inreplace "pkg/scwversion/placeholder.go" do |s|
      s.gsub! /VERSION = "master"/, "VERSION = \"v#{version}\""
      s.gsub! /GITCOMMIT = "master"/, "GITCOMMIT = \"v#{version}\""
    end
    system "go", "build", "-o", "scw", "./cmd/scw"
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
