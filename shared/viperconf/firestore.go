package viperconf

import (
	"context"
	"log"
	"os"
	"strings"

	"cloud.google.com/go/firestore"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

const secretPrefixFirestore = "FIRESTORE:"

func GetFirestoreFullConfig(thisViper *viper.Viper, cfgLocation string, setDefault bool) error {
	remoteVals := strings.SplitN(cfgLocation, "|", 2)
	if len(remoteVals) < 2 {
		log.Printf("Invalid firestore REMOTE_CONFIG format. Must be in format \"" + secretPrefixFirestore + "product-id|path\"")
		os.Exit(1)
	}
	projectID := remoteVals[0]
	path := remoteVals[1]

	firestoreClient, err := firestore.NewClient(context.Background(), projectID) // connectionURI == projectID
	if err != nil {
		log.Printf("NewClient error: %v", err)
		os.Exit(1)
	}

	fsobj, err := firestoreClient.Doc(path).Get(context.Background())
	if err != nil {
		if grpc.Code(err) == codes.NotFound {
			// log.Printf("Got 404 error: %v", err)
			return err
		}
		log.Printf("Get error: %v", err)
		return err
	}

	tmpmap := make(map[string]interface{})
	err = fsobj.DataTo(&tmpmap)
	if err != nil {
		log.Printf("firestore DataTo ERROR: %v", err)
		return err
	}

	for k, v := range tmpmap {
		if setDefault {
			thisViper.SetDefault(k, v)
		} else {
			thisViper.Set(k, v)
		}
	}

	return nil
}
