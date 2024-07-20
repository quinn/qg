class G < Formula
    desc "Description of your Go binary"
    homepage "https://github.com/quinn/g"
    version "0.1.0"
  
    if OS.mac?
      url "https://github.com/quinn/g/releases/download/v0.1.0/g-darwin-amd64"
      sha256 "SHA256_SUM_OF_DARWIN_BINARY"
    elsif OS.linux?
      url "https://github.com/quinn/g/releases/download/v0.1.0/g-linux-amd64"
      sha256 "SHA256_SUM_OF_LINUX_BINARY"
    end
  
    def install
      if OS.mac?
        bin.install "g-darwin-amd64" => "g"
      elsif OS.linux?
        bin.install "g-linux-amd64" => "g"
      end
    end
  
    test do
      system "#{bin}/g", "--version"
    end
  end
  