import sys
# insert path of flask_dev dir 
sys.path.insert(0, '/var/www/flask_dev')
# import our app from file server
from server import app as application

import sys
sys.stdout = sys.stderr
