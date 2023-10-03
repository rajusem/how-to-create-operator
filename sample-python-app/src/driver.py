import datetime
import platform

def main():
    print('inside main method')
    print(f"Current UTC time: {datetime.datetime.now()}")
    print(f"Platform info : {platform.node()}")

if __name__ == '__main__':
    main()
