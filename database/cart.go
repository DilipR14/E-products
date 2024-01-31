package database

import (
	"context"
	"errors"
	"log"
	"time"

	
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	ErrCanFindProduct        = errors.New("can't find the product")
	ErrCanDecodeProducts     = errors.New("can't find the product")
	ErrUserIdNotValid        = errors.New("this user is not valid")
	ErrCanUpdateUser         = errors.New("cannot add this product to the cart")
	ErrCanRemoveItemCart     = errors.New("cannot remove this item from the cart")
	ErrCanGetItem            = errors.New("was unable to get the item from the cart")
	ErrCanBuyCartItem        = errors.New("Cannot update the purchase")
)

func AddProductToCart(ctx context.Context, prodCollection, UserCollection *mongo.Collection, productID primitive.ObjectID, userID string) error {
	searchfromdb, err := prodCollection.Find(ctx, bson.M{"_id": productID})
	if err != nil {
		log.Println(err)
		return ErrCanFindProduct
	}
	var productCart []models.ProductUser

	err = searchfromdb.All(ctx, &productCart)
	if err != nil {
		log.Println(err)
		return ErrCanDecodeProducts
	}
	id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		log.Println(err)
		return ErrUserIdNotValid
	}

	filter := bson.D{primitive.E{Key: "_id", Value: id}}
	update := bson.D{{Key: "$push", Value: bson.D{primitive.E{Key: "usercart", Value: bson.D{{Key: "$each", Value: productCart}}}}}}

	_, err = UserCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return ErrUserIdNotValid
	}
	return nil
}

func RemoveItemToCart(ctx context.Context, prodCollection, UserCollection *mongo.Collection, productID primitive.ObjectID, userID string) error {

	id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		log.Println(err)
		return ErrUserIdNotValid
	}
	filter := bson.D{primitive.E{Key: "_id", Value: id}}
	update := bson.D{"$pull": bson.M{"usercart": bson.M{"_id": productID}}}

	_, err := UserCollection.UpdateMany(ctx, filter, update)
	if err != nil {
		return ErrCanRemoveItemCart
	}
	return nil
}

func BuyItemFromCart(ctx context.Context, UserCollection *mongo.Collection, userID string) error {

	id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		log.Println(err)
		return ErrUserIdNotValid
	}

	var getcartitems models.User
	var ordercart models.Order

	ordercart.Order_ID = primitive.NewObjectID()
	ordercart.Order_At = time.Now()
	ordercart.Order_Cart = make([]models.ProductUser, 0)
	ordercart.PaymentMethod.COD = true

	unwind := bson.D{{Key: "$unwind", Value: bson.D{primitive.E{Key: "path", Value: "$usercart"}}}}
	grouping := bson.D{{Key: "$group", Value: bson.D{primitive.E{Key: "_id", Value: "$_id"}, {Key: "total", Value: bson.D{primitive.E{Key: "$sum", Value: "$usercart.price"}}}}}}

	currentresults, err := UserCollection.Aggregate(ctx, mongo.Pipeline{unwind, grouping})
	ctx.Done()
	if err != nil {
		panic(err)
	}

	var getusercart []bson.M
	currentresults.All(ctx, &getusercart)
	if err != nil {
		panic(err)
	}

	var total_price int32

	for _, userItem := range getusercart {
		price := userItem["total"]
		total_price = price.(int32)
	}
	ordercart.Price = int(total_price)

	filter := bson.D{primitive.E{Key: "_id", Value: id}}
	update := bson.D{{Key: "$push", Value: bson.D{primitive.E{Key: "orders", Value: ordercart}}}}

	_, err = UserCollection.UpdateMany(ctx, filter, update)

	if err != nil {
		log.Println(err)
	}

	err = UserCollection.FindOne(ctx, bson.D{primitive.E{Key: "_id", Value: id}}).Decode(&getcartitems)
	if err != nil {
		log.Println(err)
	}
	filter2 := bson.D{primitive.E{Key: "_id", Value: id}}
	update2 := bson.M{"$push": bson.M{"order.$[].order_list": bson.M{"$each": getcartitems.UserCart}}}
	_, err = UserCollection.UpdateOne(ctx, filter2, update2)

	if err != nil {
		log.Println(err)
	}

	usercartEmpty := make([]models.ProductUser, 0)
	filter3 := bson.D{primitive.E{Key: "_id", Value: id}}
	update3 := bson.D{{Key: "$set", Value: bson.D{primitive.E{Key: "usercart", Value: usercartEmpty}}}}
	_, err = UserCollection.UpdateOne(ctx, filter3, update3)
	if err != nil {
		return ErrCanBuyCartItem
	}

	return nil
}

func InstantBuyer(ctx context.Context, prodCollection, UserCollection *mongo.Collection, productID primitive.ObjectID, userID string) error {

	id, err := primitive.ObjectIDFromHex(userID)

	if err != nil {
		log.Println(err)
		return ErrUserIdNotValid
	}
	var productDetails models.ProductUser
	var orderDetails models.Order

	orderDetails.Order_ID = primitive.NewObjectID()
	orderDetails.Order_At = time.Now()
	orderDetails.Order_Cart = make([]models.ProductUser, 0)
	orderDetails.PaymentMethod.COD = true
	prodCollection.FindOne(ctx, bson.D{primitive.E{Key: "_id", Value: productID}}).Decode(&productDetails)

	if err != nil {
		log.Println(err)
	}
	orderDetails.Price = productDetails.Price

	filter := bson.D{primitive.E{Key: "_id", Value: id}}
	update := bson.D{{Key: "$push", Value: bson.D{primitive.E{Key: "orders", Value: orderDetails}}}}
	_, err = UserCollection.UpdateMany(ctx, filter, update)
	if err != nil {
		log.Println(err)
	}

	filter2 := bson.D{primitive.E{Key: "_id", Value: id}}
	update2 := bson.M{"$push": bson.M{"orders.$[].order_list": productDetails}}
	_, err = UserCollection.UpdateMany(ctx, filter2, update2)
	if err != nil {
		log.Println(err)
	}
	return nil
}
