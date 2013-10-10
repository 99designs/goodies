
module Sprockets
  class StaticAsset < Asset
    alias :write_to_with_compression :write_to
    def write_to(filename, options = {})
      return if File.extname(filename) == '.gz'
      write_to_with_compression(filename, options)
    end
  end
end
