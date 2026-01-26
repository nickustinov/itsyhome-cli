class Itsyhome < Formula
  desc "CLI tool to control HomeKit devices via Itsyhome"
  homepage "https://github.com/nickustinov/itsyhome-cli"
  url "https://github.com/nickustinov/itsyhome-cli/archive/refs/tags/v0.3.0.tar.gz"
  sha256 "e0a853e65e8f9d4254bce331022da9e1d3961253da1313dc753acc71dbc6e345"
  license "MIT"

  depends_on "go" => :build

  def install
    system "go", "build", *std_go_args(ldflags: "-s -w -X main.version=#{version}")
  end

  test do
    assert_match "itsyhome", shell_output("#{bin}/itsyhome --help")
  end
end
