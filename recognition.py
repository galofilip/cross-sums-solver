import sys
import pytesseract
from PIL import Image

def main():
	arr = ""
	custom_config = r'--psm 6 -c tessedit_char_whitelist=0123456789'
	s = sys.argv[1].rstrip(",").split(",")
	for file in s:
		try:
			img = Image.open(file)
		except Exception:
			arr += "-1,"
			continue
		text = pytesseract.image_to_string(img, config=custom_config).strip()
		arr += text+","
	print(arr)

if __name__ == "__main__":
    main()
