# coding: utf-8

lib = File.expand_path('../lib', __FILE__)
$LOAD_PATH.unshift(lib) unless $LOAD_PATH.include?(lib)

Gem::Specification.new do |spec|
  spec.name = "go_lang_qr_code".freeze
  spec.version = "1.0.0"
  spec.authors = ["GARAIO REM".freeze]
  spec.date = "2021-05-27"
  spec.description = "A QR code library that works with regular Ruby. The QR code is generated with Golang because of much better performance. This gem uses FFI as an interface between Golang and Ruby. It is forked from https://github.com/divan/qrlogo".freeze
  spec.summary = "A QR code library".freeze
  spec.homepage = "https://github.com/garaio/go_lang_qr_code".freeze
  spec.licenses = ["MIT".freeze]
  spec.platform = Gem::Platform::RUBY
  spec.required_ruby_version = '>= 2.5.0'
  spec.required_rubygems_version = '>= 2.0.0'

  spec.files = `git ls-files -z`.split("\x0").reject do |f|
    f.match(%r{^(test|spec|features)/})
  end

  spec.add_dependency("ffi", ["~> 1.13"])
end
