package main
import (
	"time"
	"context"
	"fmt"
	"html/template"
	"net/http"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"
)

type Events struct {
	Event string
	Time string
}

type Mesg struct {
        Msg string `json:"message" binding:"required"`
}

func welcome(c *gin.Context) {
	c.String(http.StatusOK, "Welcome to Devops Dashboard")
}

func get_events(c *gin.Context) {
	var msg string
	findOptions := options.Find()
	findOptions.SetLimit(100)
//	var events []*Events
        clientOptions := options.Client().ApplyURI("mongodb://admin:admin123@54.205.212.185:27017/admin")
//      Connect to MongoDB
        client, err := mongo.Connect(context.TODO(), clientOptions)
        if err != nil {
		fmt.Println(err)
        }
        collection := client.Database("event_dash").Collection("events")
// 	Passing bson.D{{}} as the filter matches all documents in the collection
	cur, err := collection.Find(context.TODO(), bson.M{}, findOptions)
	if err != nil {
    		fmt.Println(err)
	}
// 	Finding multiple documents returns a cursor
//	Iterating through the cursor allows us to decode documents one at a time
	for cur.Next(context.TODO()) {
    
// 	create a value into which the single document can be decoded
    	var elem Events
    	err := cur.Decode(&elem)
    	if err != nil {
        		fmt.Println(err)
    		}
//    	events = append(events, &elem)
	msg = fmt.Sprintf("%s<tr><td> <strong class=\"event-title\">%s</strong></td><td class=\"event-time\">%s</td></tr>",msg,elem.Event,elem.Time)
	}
	if err := cur.Err(); err != nil {
    		fmt.Println(err)
	}
// 	Close the cursor once finished
	defer cur.Close(context.TODO())
//	e := "Event 1 has come"
//	t := "00:00 UTC"
//	msg = fmt.Sprintf("<tr><td> <strong class=\"event-title\">%s</strong></td><td class=\"event-time\">%s</td></tr>",e,t)
	cont := template.HTML(msg)
	c.HTML(http.StatusOK, "temp.tmpl", gin.H{"val":cont,})	
}

func get_details(c *gin.Context) {
	eid := c.DefaultQuery("id","nil")
	c.String(http.StatusOK, "id = %s", eid)
}

func put_event(c *gin.Context) {

	loc, err := time.LoadLocation("Asia/Kolkata")
	if err != nil {
        	fmt.Println(err)
    	}
	et := time.Now().In(loc).String()
	var json Mesg
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

//	Mongo DB Connect
//	Set client options
	clientOptions := options.Client().ApplyURI("mongodb://admin:admin123@54.205.212.185:27017/admin")

//	Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
    		
	}

	collection := client.Database("event_dash").Collection("events")
//	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)

	e1 := Events{json.Msg,et}
	insertResult, err := collection.InsertOne(context.TODO(), e1)
	if err != nil {
    	
	}
	fmt.Println("Inserted an event: ", json.Msg, " Ref: ", insertResult.InsertedID)
	c.JSON (200, gin.H{
			"status": "posted",
			"message": json.Msg,
	})
}

func main() {
	
	router := gin.New()
	gin.ForceConsoleColor()	
	router.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {

		// your custom format
		return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
				param.ClientIP,
				param.TimeStamp.Format(time.RFC1123),
				param.Method,
				param.Path,
				param.Request.Proto,
				param.StatusCode,
				param.Latency,
				param.Request.UserAgent(),
				param.ErrorMessage,
		)
	}))
	router.Use(gin.Recovery())
//	router.HandleFunc("/", welcome).Methods("GET") ## for MUX

//	Mongo DB Connect
//	Set client options
	clientOptions := options.Client().ApplyURI("mongodb://admin:admin123@54.205.212.185:27017/admin")

//	Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
    		
	}

//	Check the connection
	err = client.Ping(context.TODO(), nil)

	if err != nil {
    		
	}

	fmt.Println("Connected to MongoDB!")

	router.LoadHTMLGlob("template/*")
//	router.LoadHTMLFiles("templates/template1.html","templates/template2.html")
	router.Static("/css","./template")

	router.GET("/", welcome)
	router.GET("/api/get_events", get_events)
	router.GET("/api/get_details", get_details)
	router.POST("/api/put_event", put_event)

	router.Run(":5000")
	fmt.Println("Listening on port 5000")
}
