package spa

import (
	"fmt"
	"github.com/gobuffalo/packr/v2"
	"github.com/iotaledger/hive.go/parameter"
	"github.com/labstack/echo"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
	"time"
)

var ErrInvalidParameter = errors.New("invalid parameter")
var ErrInternalError = errors.New("internal error")
var ErrNotFound = errors.New("not found")
var ErrForbidden = errors.New("forbidden")

// holds SPA assets
var appBox = packr.New("SPA_App", "./frontend/build")
var assetsBox = packr.New("SPA_Assets", "./frontend/src/assets")

func indexRoute(e echo.Context) error {
	if parameter.NodeConfig.GetBool("dashboard.dev") {
		res, err := http.Get("http://127.0.0.1:9090/")
		if err != nil {
			return err
		}
		devIndexHTML, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return err
		}
		return e.HTMLBlob(http.StatusOK, devIndexHTML)
	}
	indexHTML, err := appBox.Find("index.html")
	if err != nil {
		return err
	}
	return e.HTMLBlob(http.StatusOK, indexHTML)
}

func setupRoutes(e *echo.Echo) {

	if parameter.NodeConfig.GetBool("dashboard.dev") {
		e.Static("/assets", "./plugins/spa/frontend/src/assets")
	} else {

		// load assets from packr: either from within the binary or actual disk
		e.GET("/app/*", echo.WrapHandler(http.StripPrefix("/app", http.FileServer(appBox))))
		e.GET("/assets/*", echo.WrapHandler(http.StripPrefix("/assets", http.FileServer(assetsBox))))
	}

	e.GET("/ws", websocketRoute)
	e.GET("/", indexRoute)

	// used to route into the SPA index
	e.GET("*", indexRoute)

	apiRoutes := e.Group("/api")

	setupExplorerRoutes(apiRoutes)

	e.HTTPErrorHandler = func(err error, c echo.Context) {
		c.Logger().Error(err)

		var statusCode int
		var message string

		switch errors.Cause(err) {

		case echo.ErrNotFound:
			c.Redirect(http.StatusSeeOther, "/")
			return

		case echo.ErrUnauthorized:
			statusCode = http.StatusUnauthorized
			message = "unauthorized"

		case ErrForbidden:
			statusCode = http.StatusForbidden
			message = "access forbidden"

		case ErrInternalError:
			statusCode = http.StatusInternalServerError
			message = "internal server error"

		case ErrNotFound:
			statusCode = http.StatusNotFound
			message = "not found"

		case ErrInvalidParameter:
			statusCode = http.StatusBadRequest
			message = "bad request"

		default:
			statusCode = http.StatusInternalServerError
			message = "internal server error"
		}

		message = fmt.Sprintf("%s, error: %+v", message, err)
		c.String(statusCode, message)
	}
}

func registerWSClient() (uint64, chan interface{}) {
	// allocate new client id
	clientsMu.Lock()
	defer clientsMu.Unlock()
	clientID := nextClientID
	channel := make(chan interface{}, 50)
	preFeed(channel)
	clients[clientID] = channel
	nextClientID++
	return clientID, channel
}

func websocketRoute(c echo.Context) error {
	defer func() {
		if r := recover(); r != nil {
			log.Errorf("recovered from panic within WS handle func: %s", r)
		}
	}()
	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}
	defer ws.Close()

	// cleanup client websocket
	clientID, channel := registerWSClient()
	defer func() {
		clientsMu.Lock()
		delete(clients, clientID)
		close(channel)
		clientsMu.Unlock()
	}()

	for {
		if err := ws.WriteJSON(<-channel); err != nil {
			log.Errorf("error while writing to web socket client %s: %s", c.RealIP(), err.Error())
			break
		}
		if err := ws.SetWriteDeadline(time.Now().Add(webSocketWriteTimeout)); err != nil {
			log.Errorf("error while setting write deadline on web socket client %s: %s", c.RealIP(), err.Error())
			break
		}
	}
	return nil
}
