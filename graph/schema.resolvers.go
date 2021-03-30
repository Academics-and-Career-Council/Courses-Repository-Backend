package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"log"

	"github.com/Mshivam2409/AnC-Courses/graph/generated"
	"github.com/Mshivam2409/AnC-Courses/models"
	"github.com/Mshivam2409/AnC-Courses/services"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
)

func (r *mutationResolver) AddCourse(ctx context.Context, course models.NewCourse) (*models.Response, error) {
	// i := fmt.Sprintf("%v", ctx.Value("Authorizaton"))
	// a, err := services.OryClient.IsAuthorized(i, "read", "course")
	// if !a {
	// 	return &models.Response{Ok: false, Message: "Unauthorzed"}, err
	// }
	d := services.GetService()
	fid, err := d.CreateFolder(course.Number)
	log.Println(fid)
	if err != nil {
		log.Println(err, fid)
		return &models.Response{}, err
	}
	c := &models.MGMCourse{Title: course.Title,
		Author:      course.Author,
		Credits:     course.Credits,
		Contents:    course.Contents,
		Offered:     course.Offered,
		Dept:        course.Dept,
		DriveFolder: fid,
		DriveFiles:  []string{},
		Number:      course.Number,
		Reviews:     []primitive.ObjectID{},
	}
	_, err = services.MongoClient.Courses.Collection("courses").InsertOne(context.TODO(), c)
	return &models.Response{Ok: true, Message: "Course created Successfully!"}, err
}

func (r *mutationResolver) AddReview(ctx context.Context, review models.NewReview) (*models.Response, error) {
	wc := writeconcern.New(writeconcern.WMajority())
	rc := readconcern.Snapshot()
	txnOpts := options.Transaction().SetWriteConcern(wc).SetReadConcern(rc)
	session, err := services.MongoClient.Courses.Client().StartSession()
	if err != nil {
		return &models.Response{Ok: false, Message: err.Error()}, err
	}
	defer session.EndSession(context.Background())
	err = mongo.WithSession(context.Background(), session, func(sessionContext mongo.SessionContext) error {
		c := &models.MGMCourse{}
		filter := bson.D{{Key: "number", Value: review.Course}}
		err := services.MongoClient.Courses.Collection("courses").FindOne(sessionContext, filter).Decode(c)
		if err != nil {
			log.Println(err)
			return err
		}
		cid, err := primitive.ObjectIDFromHex(c.ID)
		fmt.Println(cid)
		if err = session.StartTransaction(txnOpts); err != nil {
			log.Println(err)
			return err
		}
		result, err := services.MongoClient.Courses.Collection("reviews").InsertOne(
			sessionContext,
			&models.MGMReview{
				Course:    cid,
				Contents:  review.Grading,
				Professor: review.Instructor,
				Approved:  false,
			},
		)
		fmt.Println(result.InsertedID)
		if err != nil {
			log.Println(err)
			return err
		}
		// rid := fmt.Sprintf("%v", result.InsertedID)
		fmt.Println(fmt.Sprintf("%v", result.InsertedID))
		if err != nil {
			log.Println(err)
			return err
		}
		update := bson.M{
			"$push": bson.M{
				"reviews": result.InsertedID,
			},
		}
		fmt.Printf("%v", update)
		_ = services.MongoClient.Courses.Collection("courses").FindOneAndUpdate(sessionContext, filter, update)
		fmt.Printf("1000")
		if err = session.CommitTransaction(sessionContext); err != nil {
			log.Println(err)
			return err
		}
		fmt.Println(result.InsertedID)
		return nil
	})
	if err != nil {
		if abortErr := session.AbortTransaction(context.Background()); abortErr != nil {
			return &models.Response{Ok: false, Message: err.Error()}, err
		}
		return &models.Response{Ok: false, Message: err.Error()}, err
	}
	return &models.Response{Ok: true, Message: "success"}, err
}

func (r *mutationResolver) ModifyCourse(ctx context.Context, course *models.ModifyCourseInput) (*models.Course, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ModifyReview(ctx context.Context, reviewID string, status bool) (*models.Review, error) {
	rid, err := primitive.ObjectIDFromHex(reviewID)
	rev := &models.MGMReview{}
	res := &models.Review{}
	if err != nil {
		return res, err
	}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "approved", Value: status}}}}
	err = services.MongoClient.Courses.Collection("reviews").FindOneAndUpdate(context.TODO(),
		bson.D{{Key: "_id", Value: rid}}, update).Decode(rev)
	// err = services.MongoClient.Courses.Collection("reviews").FindOne(context.TODO(), bson.D{{Key: "_id", Value: rid}}).Decode(rev)
	if err != nil {
		return res, err
	}
	res.Approved = status
	res.Course = rev.Course.Hex()
	res.Grading = rev.Contents
	res.ID = reviewID
	res.Instructor = rev.Professor
	res.Semester = rev.Semester
	return res, err
}

func (r *mutationResolver) ElevateUser(ctx context.Context, username string) (*models.User, error) {
	u := &models.MGMUser{}
	filter := bson.D{{Key: "username", Value: username}}
	update := bson.D{{Key: "$inc", Value: bson.D{{Key: "clearance", Value: 1}}}}
	err := services.MongoClient.Users.Collection("ug").FindOneAndUpdate(context.TODO(), filter, update).Decode(u)
	return &models.User{
		ID:        u.ID,
		Name:      u.Name,
		Rollno:    u.RollNo,
		Username:  u.Username,
		Banned:    u.Banned,
		Clearance: u.Clearance,
	}, err
}

func (r *mutationResolver) DemoteUser(ctx context.Context, username string) (*models.User, error) {
	u := &models.MGMUser{}
	filter := bson.D{{Key: "username", Value: username}}
	update := bson.D{{Key: "$inc", Value: bson.D{{Key: "clearance", Value: -1}}}}
	err := services.MongoClient.Users.Collection("ug").FindOneAndUpdate(context.TODO(), filter, update).Decode(u)
	return &models.User{
		ID:        u.ID,
		Name:      u.Name,
		Rollno:    u.RollNo,
		Username:  u.Username,
		Banned:    u.Banned,
		Clearance: u.Clearance,
	}, err
}

func (r *mutationResolver) ToggleBanUser(ctx context.Context, username string, banned bool) (*models.User, error) {
	u := &models.MGMUser{}
	filter := bson.D{{Key: "username", Value: username}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "banned", Value: banned}}}}
	err := services.MongoClient.Users.Collection("ug").FindOneAndUpdate(context.TODO(), filter, update).Decode(u)
	return &models.User{
		ID:        u.ID,
		Name:      u.Name,
		Rollno:    u.RollNo,
		Username:  u.Username,
		Banned:    u.Banned,
		Clearance: u.Clearance,
	}, err
}

func (r *queryResolver) GetCourseData(ctx context.Context, number string) (*models.CourseData, error) {
	// i := fmt.Sprintf("%v", ctx.Value("Authorizaton"))
	// a, err := services.OryClient.IsAuthorized(i, "read", "course")
	// if !a {
	// 	return &models.CourseData{}, err
	// }
	c := &models.MGMCourse{}
	filter := bson.D{{Key: "number", Value: number}}
	err := services.MongoClient.Courses.Collection("courses").FindOne(context.TODO(), filter).Decode(c)
	if err != nil {
		log.Println(err)
		return &models.CourseData{}, err
	}
	cid, err := primitive.ObjectIDFromHex(c.ID)
	if err != nil {
		log.Println(err)
		return &models.CourseData{}, err
	}
	filter = bson.D{{Key: "course", Value: cid}}
	cur, err := services.MongoClient.Courses.Collection("reviews").Find(context.TODO(), filter)
	if err != nil {
		log.Println(err)
		return &models.CourseData{}, err
	}
	reviews := []*models.Review{}
	for cur.Next(context.TODO()) {
		elem := &models.MGMReview{}
		err := cur.Decode(elem)
		if err != nil {
			log.Println(err)
			return &models.CourseData{}, err
		}
		reviews = append(reviews, &models.Review{ID: elem.ID, Semester: elem.Semester, Instructor: elem.Professor, Grading: elem.Contents, Approved: elem.Approved})
	}
	if err := cur.Err(); err != nil {
		if err != nil {
			log.Println(err)
			return &models.CourseData{}, err
		}
	}
	cur.Close(context.TODO())
	d := &models.CourseData{
		Course: &models.Course{
			Title:      c.Title,
			Number:     c.Number,
			Credits:    c.Credits,
			Offered:    c.Offered,
			Contents:   c.Contents,
			DriveFiles: c.DriveFiles,
			ID:         c.ID,
			Dept:       c.Dept,
			Author:     c.Author,
		},
		Reviews: reviews,
	}
	return d, err
}

func (r *queryResolver) GetReviewsbyCourse(ctx context.Context, number string) ([]*models.Review, error) {
	c := &models.MGMCourse{}
	filter := bson.D{{Key: "number", Value: number}}
	err := services.MongoClient.Courses.Collection("courses").FindOne(context.TODO(), filter).Decode(c)
	if err != nil {
		print(err.Error())
	}
	cid, err := primitive.ObjectIDFromHex(c.ID)
	filter = bson.D{{Key: "course", Value: cid}}
	cur, err := services.MongoClient.Courses.Collection("reviews").Find(context.TODO(), filter)
	reviews := []*models.Review{}
	for cur.Next(context.TODO()) {
		elem := &models.MGMReview{}
		err := cur.Decode(elem)
		if err != nil {
			log.Fatal(err)
		}
		reviews = append(reviews, &models.Review{ID: elem.ID, Semester: elem.Semester})
	}

	if err := cur.Err(); err != nil {
		print(err.Error())
	}
	cur.Close(context.TODO())
	return reviews, err
}

func (r *queryResolver) SearchCourses(ctx context.Context, params *models.SearchParams) ([]*models.Course, error) {
	courses := []*models.Course{}
	if params.Identifier == "" {
		return courses, nil
	}
	filter := bson.D{{Key: "number", Value: bson.D{{Key: "$regex", Value: primitive.Regex{Pattern: "^" + params.Identifier, Options: "i"}}}}}
	limit := int64(10)
	cur, err := services.MongoClient.Courses.Collection("courses").Find(context.TODO(), filter, &options.FindOptions{Limit: &limit})
	if err != nil {
		print(err.Error())
	}
	for cur.Next(context.TODO()) {
		elem := &models.MGMCourse{}
		err := cur.Decode(elem)
		if err != nil {
			log.Fatal(err)
		}

		courses = append(courses, &models.Course{ID: elem.ID, Number: elem.Number, Credits: elem.Credits, Title: elem.Title})
	}
	if err := cur.Err(); err != nil {
		print(err.Error())
	}
	cur.Close(context.TODO())
	return courses, err
}

func (r *queryResolver) SearchUsers(ctx context.Context, params *models.SearchParams) ([]*models.User, error) {
	panic(fmt.Errorf("not implemented"))
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
