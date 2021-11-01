# run it with
# irb
# require_relative 'lib/go_lang/qr_code_test'
# GoLang::QrCodeTest.test
require_relative 'qr_code'

module GoLang::QrCodeTest
	def self.test
		qr_string = <<-data
			SPC
			0200
			1
			CH4431999123000889012
			K
			STWEG Grenzweg
			-
			5036 Oberentfelden


			CH







			854.00
			CHF
			K
			Fimospa AG
			Kantonsstrasse
			6234 Triengen


			CH
			QRR
			834382000000040000000078358

			EPD
			data

			bild = "/home/schaerli/Dev/rem2/config/qr_code/ch-kreuz_7mm/ch-kreuz_go.png"

			1000.times do
				puts "*" * 100
				puts `ps -o rss= -p #{$$}`.to_i / 1000
				puts "*" * 100
				GoLang::QrCode.create_qr_code_with_logo(qr_string, overlay_logo_path: bild, qr_code_size: 512)
			end
			# in einem geforkten Prozess bleibt go in einem infity loop hängen im garbage colletor
			# fork do
			# 	1000.times do
			# 		puts "*" * 100
			# 		puts `ps -o rss= -p #{$$}`.to_i / 1000
			# 		puts "*" * 100
			# 		GoLang::QrCode.create_qr_code_with_logo(qr_string, overlay_logo_path: bild, qr_code_size: 512)
			# 	end
			# end
	end
end