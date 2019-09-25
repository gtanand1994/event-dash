package main
import (
        "context"
        "fmt"
        "go.mongodb.org/mongo-driver/mongo"
        "go.mongodb.org/mongo-driver/mongo/options"
        "gopkg.in/mgo.v2/bson"
)

type Events struct {
	Event string
	Time string
}

func main() {
	var msg string
	findOptions := options.Find()
	findOptions.SetLimit(2)
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
	fmt.Println(msg)
}
