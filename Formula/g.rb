class Mygobin < Formula
    desc "Description of your Go binary"
    homepage "https://github.com/quinn/mygobin"
    version "0.1.0"
  
    if OS.mac?
      url "https://github.com/quinn/mygobin/releases/download/v0.1.0/mygobin-darwin-amd64"
      sha256 "SHA256_SUM_OF_DARWIN_BINARY"
    elsif OS.linux?
      url "https://github.com/quinn/mygobin/releases/download/v0.1.0/mygobin-linux-amd64"
      sha256 "SHA256_SUM_OF_LINUX_BINARY"
    end
  
    def install
      if OS.mac?
        bin.install "mygobin-darwin-amd64" => "mygobin"
      elsif OS.linux?
        bin.install "mygobin-linux-amd64" => "mygobin"
      end
    end
  
    test do
      system "#{bin}/mygobin", "--version"
    end
  end
  