require "language/go"

class Anonuuid < Formula
  desc "Anonymize streams"
  homepage "https://github.com/moul/anonuuid"
  url "https://github.com/moul/anonuuid/archive/v1.0.0.tar.gz"
  sha256 "37d6ff3931276c7e8eac6d8d5c34d0deb1a649567dc2c4756fc4f74637a0eb51"

  head "https://github.com/moul/anonuuid.git"

  depends_on "go" => :build

  def install
    ENV["GOPATH"] = buildpath
    ENV["CGO_ENABLED"] = "0"
    ENV.prepend_create_path "PATH", buildpath/"bin"

    mkdir_p buildpath/"src/github.com/moul"
    ln_s buildpath, buildpath/"src/github.com/moul/anonuuid"
    Language::Go.stage_deps resources, buildpath/"src"

    system "go", "get", "github.com/codegangsta/cli"
    system "go", "build", "-o", "anonuuid", "./cmd/anonuuid"
    bin.install "anonuuid"

    # FIXME: add autocompletion
  end

  test do
    output = shell_output(bin/"anonuuid --version")
    assert output.include? "anonuuid version"
  end
end
