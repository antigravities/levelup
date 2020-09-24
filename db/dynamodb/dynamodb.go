package dynamodb

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"get.cutie.cafe/levelup/types"
	"get.cutie.cafe/levelup/util"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

var (
	sess *session.Session
	db   *dynamodb.DynamoDB

	// Cache stores apps obtained this session
	Cache map[int]types.App = make(map[int]types.App)

	table string
)

// Initialize and connect to DynamoDB
func Initialize() {
	util.Info("Initializing DynamoDB")
	sess, err := session.NewSession()
	if err != nil {
		panic(err)
	}

	table = os.Getenv("LU_TABLE")

	db = dynamodb.New(sess)
}

// GetApp gets the information about an app.
func GetApp(appid int) *types.App {
	util.Info(fmt.Sprintf("Fetching app %d", appid))

	if val, ok := Cache[appid]; ok {
		util.Debug(fmt.Sprintf("Cache: hit"))
		return &val
	}

	util.Debug(fmt.Sprintf("Cache: miss"))

	result, err := db.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(table),
		Key: map[string]*dynamodb.AttributeValue{
			"AppID": {
				N: aws.String(strconv.Itoa(appid)),
			},
		},
	})

	if err != nil {
		return nil
	}

	app := types.App{}

	if err := dynamodbattribute.UnmarshalMap(result.Item, &app); err != nil {
		return nil
	}

	Cache[appid] = app

	return &app
}

// GetApps returns all of the AppIDs in the table.
func GetApps(pending bool) []int {
	util.Info("Fetching apps from DynamoDB")

	util.Debug(fmt.Sprintf("pending: %v", pending))

	input := &dynamodb.ScanInput{
		TableName:            aws.String(table),
		ProjectionExpression: aws.String("AppID"),
		FilterExpression:     aws.String("IsPending <> :f"),
	}

	if !pending {
		input.ExpressionAttributeValues = map[string]*dynamodb.AttributeValue{
			":f": {
				BOOL: aws.Bool(true),
			},
		}
	} else {
		input.ExpressionAttributeValues = map[string]*dynamodb.AttributeValue{
			":f": {
				BOOL: aws.Bool(false),
			},
		}
	}

	res, err := db.Scan(input)

	if err != nil {
		util.Warn(fmt.Sprintf("Error: %v", err))
		return []int{}
	}

	apps := []int{}

	for _, item := range res.Items {
		appid := 0

		if err := dynamodbattribute.Unmarshal(item["AppID"], &appid); err != nil {
			continue
		}

		apps = append(apps, appid)
	}

	return apps
}

// GetCachedApps filters through the cache for bad apps and returns only good ones.
// TODO: figure out how bad apps get in the cache anyway
func GetCachedApps() map[int]types.App {
	apps := make(map[int]types.App)

	for _, v := range Cache {
		if v.IsPending || v.AppID == 0 || v.Name == "" {
			continue
		}

		apps[v.AppID] = v
	}

	return apps
}

// PutApp updates or creates an app in the table with new information from a *types.App.
func PutApp(app types.App) error {
	util.Info(fmt.Sprintf("Putting app %d", app.AppID))

	if app.AppID == 0 {
		util.Warn("Trying to store an app with ID 0, cancelling")
		return nil
	}

	av, err := dynamodbattribute.MarshalMap(app)
	if err != nil {
		return err
	}

	if _, err = db.PutItem(&dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(table),
	}); err != nil {
		util.Warn(fmt.Sprintf("Error putting app: %v", err))
		return err
	}

	if !app.IsPending || app.AppID == 0 || app.Name == "" {
		Cache[app.AppID] = app
		util.Debug("Cache: stored")
	}

	return nil
}

// FindStaleApps finds apps that haven't been updated in an hour or more
func FindStaleApps() []types.App {
	util.Info("Finding stale apps")

	res, err := db.Scan(&dynamodb.ScanInput{
		TableName:            aws.String(table),
		ProjectionExpression: aws.String("AppID"),
		FilterExpression:     aws.String("LastUpdate < :t AND IsPending <> :f"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":t": {
				N: aws.String(strconv.FormatInt(time.Now().Unix()-60*60, 10)),
			},
			":f": {
				BOOL: aws.Bool(true),
			},
		},
	})

	if err != nil {
		util.Warn(fmt.Sprintf("Error: %v", err))
		return []types.App{}
	}

	apps := []types.App{}

	for _, item := range res.Items {
		app := types.App{}

		if err := dynamodbattribute.UnmarshalMap(item, &app); err != nil {
			continue
		}

		apps = append(apps, app)
	}

	util.Info(fmt.Sprintf("Found %d stale apps", len(apps)))

	return apps
}

// DeleteApp deletes an app from the database.
func DeleteApp(appid int) error {
	util.Warn(fmt.Sprintf("Deleting app %d", appid))

	_, err := db.DeleteItem(&dynamodb.DeleteItemInput{
		TableName: aws.String(table),
		Key: map[string]*dynamodb.AttributeValue{
			"AppID": {
				N: aws.String(strconv.Itoa(appid)),
			},
		},
	})

	if err != nil {
		util.Warn(fmt.Sprintf("Error: %v", err))
		return err
	}

	util.Debug("Deleting app from cache")
	delete(Cache, appid)

	return nil
}
