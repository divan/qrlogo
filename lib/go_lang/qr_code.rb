require 'ffi'
require 'tempfile'

module GoLang
	module QrCode
		extend FFI::Library

		if RUBY_PLATFORM =~ /darwin/
			ffi_lib File.dirname(__FILE__) + '/qrlogo_mac.so'
		elsif RUBY_PLATFORM =~ /linux/
			ffi_lib "go_src/qrlogo.so"
		else
			raise "Plattform wird nicht unterstützt"
		end

		# define class GoString to map:
		# C type struct { const char *p; GoInt n; }
		class GoString < FFI::Struct
			layout :p,   :pointer,
						:len, :long_long
		end

		# Params: qrCodeString | qrCodePath | qrCodeSize (in Px)
		attach_function :CreateQrCode, [GoString.by_value, GoString.by_value, :long_long], :void

		# Params: qrCodeString | qrCodeSize (in Px)
		attach_function :CreateQrCodeAsBase64String, [GoString.by_value, :long_long], :strptr

		# Params: qrCodeString | qrCodePath | overlayLogoPath | qrCodeSize (in Px)
		attach_function :CreateQrCodeWithLogo, [GoString.by_value, GoString.by_value, GoString.by_value, :long_long], :void

		# Params: qrCodeString | overlayLogoPath | qrCodeSize (in Px)
		attach_function :CreateQrCodeWithLogoAsBase64String, [GoString.by_value, GoString.by_value, :long_long], :strptr

		# attach_function :login, [:string], :strptr
		attach_function :FreeUnsafePointer, [:pointer], :void


		def self.create_qr_code(qr_code_as_string, qr_code_size: 512, in_memory: false)
			# Umwandlung in Go-String (Go braucht einen Memory-Pointer)
			go_string_qr_code           = GoLang::QrCode.get_go_string_for(qr_code_as_string)

			if in_memory
				base_64_string, pointer = GoLang::QrCode.CreateQrCodeAsBase64String go_string_qr_code, qr_code_size
				FFI::AutoPointer.new(pointer, GoLang::QrCode.method(:FreeUnsafePointer))
				"data:image/png;base64,#{base_64_string}"
			else
				tempfile = Tempfile.new(['qr_code', '.png'])
				go_string_qr_code_path      = GoLang::QrCode.get_go_string_for(tempfile.path)
				GoLang::QrCode.CreateQrCode go_string_qr_code,
																		go_string_qr_code_path,
																		qr_code_size

				tempfile.path
			end
		end

		def self.create_qr_code_with_logo(qr_code_as_string, overlay_logo_path: Util::SwissQrCode::CH_KREUZ_7MM_PNG, qr_code_size: 512, in_memory: false)
			# Umwandlung in Go-String (Go braucht einen Memory-Pointer)
			go_string_qr_code           = GoLang::QrCode.get_go_string_for(qr_code_as_string)
			go_string_overlay_logo_path = GoLang::QrCode.get_go_string_for(overlay_logo_path)

			if in_memory
				base_64_string, pointer = GoLang::QrCode.CreateQrCodeWithLogoAsBase64String go_string_qr_code, go_string_overlay_logo_path, qr_code_size
				FFI::AutoPointer.new(pointer, GoLang::QrCode.method(:FreeUnsafePointer))

				"data:image/png;base64,#{base_64_string}"
			else
				tempfile                = ::Tempfile.new(['swiss_qr_code', '.png'])
				go_string_qr_code_path = GoLang::QrCode.get_go_string_for(tempfile.path)

				GoLang::QrCode.CreateQrCodeWithLogo go_string_qr_code,
																						go_string_qr_code_path,
																						go_string_overlay_logo_path,
																						qr_code_size

				tempfile.path
			end
		end


		private

		def self.get_go_string_for(a_string)
			go_string       = GoLang::QrCode::GoString.new
			go_string[:p]   = FFI::MemoryPointer.from_string(a_string)
			go_string[:len] = a_string.size
			go_string
		end
	end
end
