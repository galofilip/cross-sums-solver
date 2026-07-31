# cross-sums-solver

> ⚠️ **WARNING:** currently works only for android but i am not really planing on making a version that works for iphone

## to run:
### requirements:
#### (these are the versions of my languages but i dont know if it will work with older/newer versions but it will probbly work with newer)
1. go 1.26.5
2. python 3.14.6
3. adb cli
4. git cli

### stages:
1. turn on developer options on your phone
2. in developer options turn on wireless debuging 
3. press "Pair device with with pairing code" 
4. go to your terminal(i use termux) and write adb pair your ip:the port in the pop up(you should see successfully paired in the pop up)
5. now write in the terminal adb connect your ip:the port in IP address and port(its above where the pop up was)
6. now write in the terminal these commands in this one by one:
	git clone https://github.com/galofilip/cross-sums-solver.git
	cd cross-sums-solver
	python -m venv venv
	source venv/bin/activate
	pip install -r requirements.txt
	mkdir temp
7. then to run do go run .
