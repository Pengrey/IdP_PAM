import site

site.main()

import os
import qrcode
import sqlite3
import json
import requests
import time

def get_idps(username, DATABASE_PATH):
    conn = sqlite3.connect(DATABASE_PATH)
    c = conn.cursor()

    c.execute("SELECT idp FROM attributes WHERE username = ?", (username,))
    idps = c.fetchall()

    conn.close()
    return idps

def get_idp(idp, username, DATABASE_PATH):
    conn = sqlite3.connect(DATABASE_PATH)
    c = conn.cursor()

    c.execute("SELECT attributes FROM attributes WHERE username = ? AND idp = ?", (username, idp))
    idp = c.fetchall()

    conn.close()

    try:
        python_dict = json.loads(idp[0][0])
        idp = {
            str(key) if isinstance(key, unicode) else key:
            str(value) if isinstance(value, unicode) else value
            for key, value in python_dict.items()
        }
    except:
        print('\033[1;31m[!]\033[0m Invalid IdP')
        return None

    return idp

def parse_response(response, dict):
    response = response.replace(' ' , '') \
                       .replace('\n', '') \
                       .replace('}' , '') \
                       .replace('{' , '') \
                       .replace('":' , '=') \
                       .replace('"' , '') \
                       .replace(',' , '&') 

    for pair in response.split('&'):
        key, value = pair.split('=')
        dict[key] = value
    return dict

def request_device(client_id, scope, url):
    data = {
        'client_id': client_id,
        'scope': scope
    }

    response = requests.post(url, data=data)

    if response.status_code == 200:
        response_dict = parse_response(response.text, {})
        return (response_dict['device_code'], response_dict['user_code'])
    else:
        print('\033[1;31m[!]\033[0m Error requesting device code')
        return None

def generate_qr_code(url):
    qr = qrcode.QRCode(version=1, box_size=10, border=5)
    qr.add_data(url)
    qr.make(fit=True)
    qr.print_ascii()

def poll_for_token(url, arguments):
    interval = 5
    timeout = 60

    start_time = time.time()
    while True:
        response = requests.post(url, data=arguments)
        if response.status_code == 200:
            response_dict = parse_response(response.text, {})
            
            if 'access_token' in response_dict:
                return response_dict
            else:
                #timeout = timeout - interval
                #if timeout <=0:
                #    return None
                time.sleep(interval)
        else:
            return None

def oauth2(idp, qr_code):
    device_code, user_code = request_device(idp['request_arguments']['client_id'], idp['request_arguments']['scope'], idp['request_url'])
    if device_code == None:
        return False

    idp['poll_arguments']['device_code'] = device_code

    if qr_code:
        print('\033[1;32m[+]\033[0m Please scan the following QR code with your mobile device:')
        generate_qr_code(idp['user_url'])
    else:
        print('\033[1;32m[+]\033[0m Please visit the following URL in your browser: ' + idp['user_url'])
    
    print('\033[1;32m[+]\033[0m Enter the following code when prompted: ' + user_code)

    response_dict = poll_for_token(idp['poll_url'], idp['poll_arguments'])
    if response_dict == None:
        return False

    return True

def login(username):
    DATABASE_PATH = "/etc/project_1.sqlite"

    idps = get_idps(username, DATABASE_PATH)

    print('\033[1;33m[?]\033[0m Please select an IdP:')
    for i, idp in enumerate(idps):
        print('[' + str(i) + '] ' + idp[0])

    selection = raw_input('[>] ')

    try:
        selection = int(selection)
        if selection < 0 or selection >= len(idps):
            raise ValueError
    except ValueError:
        print('\033[1;31m[!]\033[0m Invalid selection')
        return False

    idp = idps[selection][0]

    idp = get_idp(idp, username, DATABASE_PATH)

    qr_code = raw_input('\033[1;33m[?]\033[0m Would you like to scan a QR code? [y/n] ')
    if qr_code == 'y':
        qr_code = True
    elif qr_code == 'n':
        qr_code = False
    else:
        print('\033[1;31m[!]\033[0m Invalid input')
        return False

    return oauth2(idp, qr_code)

def auth(username):
    if username != None:
        print("\033[1;32m[+]\033[0m Authenticating user: " + username)
    
    if login(username):
        print('\033[1;32m[+]\033[0m Login successful')
        return True
    else:
        print('\033[1;31m[!]\033[0m Login failed')
        return False

def get_user(pamh):
    try:
        return pamh.get_user(None)
    except pamh.exception, e:
        return e.pam_result

def pam_sm_authenticate(pamh, flags, argv):
    user = get_user(pamh)
    if user == None:
        return pamh.PAM_USER_UNKNOWN

    try:
        if auth(user) == True:
            return pamh.PAM_SUCCESS
        else:
            return pamh.PAM_AUTH_ERR
    except pamh.exception, e:
        return pamh.PAM_AUTH_ERR

def pam_sm_open_session(pamh, flags, argv):
    user = get_user(pamh)

    if user == None:
        return pamh.PAM_USER_UNKNOWN

    home_dir = pathlib.Path("/home/" + user)

    if not home_dir.exists():
        home_dir.mkdir()

    return pamh.PAM_SUCCESS

def pam_sm_close_session(pamh, flags, argv):
    return pamh.PAM_SUCCESS

def pam_sm_setcred(pamh, flags, argv):
    return pamh.PAM_SUCCESS

def pam_sm_acct_mgmt(pamh, flags, argv):  
    return pamh.PAM_SUCCESS

def pam_sm_chauthtok(pamh, flags, argv):
    return pamh.PAM_SUCCESS
