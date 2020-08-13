import dash
import dash_auth
import dash_bootstrap_components as dbc

# bootstrap theme
# https://bootswatch.com/lux/
#external_stylesheets = [dbc.themes.DARKLY]
#external_stylesheets = [dbc.themes.SLATE]
external_stylesheets = [dbc.themes.YETI]

app = dash.Dash(__name__, external_stylesheets=external_stylesheets, requests_pathname_prefix='/performance/')
#app = dash.Dash(__name__)#, external_stylesheets=external_stylesheets)
#app.css.append_css({"external_url": "https://codepen.io/chriddyp/pen/bWLwgP.css"})

server = app.server
app.config.suppress_callback_exceptions = True

#VALID_USERNAME_PASSWORD_PAIRS = {'dash_demo': 'dash@demo'}
#
#auth = dash_auth.BasicAuth(app, VALID_USERNAME_PASSWORD_PAIRS)