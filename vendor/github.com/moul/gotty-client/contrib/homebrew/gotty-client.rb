require "language/go"

class GottyClient < Formula
  desc "gotty-client: GoTTY client for your terminal"
  homepage "https://github.com/moul/gotty-client"
  url "https://github.com/moul/gotty-client/archive/v1.5.0.tar.gz"
  sha256 "fa85b4a18c6badba507604d8f592ecf452e4f3a1abaabaedc448c7f82e22d6d7"

  head "https://github.com/moul/gotty-client.git"

  depends_on "go" => :build

  def install
    ENV["GOPATH"] = buildpath
    ENV["GOBIN"] = buildpath
    ENV["GO15VENDOREXPERIMENT"] = "1"
    (buildpath/"src/github.com/moul/gotty-client").install Dir["*"]

    system "go", "build", "-o", "#{bin}/gotty-client", "github.com/moul/gotty-client/cmd/gotty-client/"

    bash_completion.install "#{buildpath}/src/github.com/moul/gotty-client/contrib/completion/bash_autocomplete"
    zsh_completion.install "#{buildpath}/src/github.com/moul/gotty-client/contrib/completion/zsh_autocomplete"
  end

  test do
    output = shell_output(bin/"gotty-client --version")
    assert output.include? "gotty-client version"
  end
end
