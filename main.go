package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	_ "goApi/docs"

	swagger "github.com/arsmn/fiber-swagger/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Location struct {
	Type        string     `json:"type" bson:"type"`
	Latitude    float64    `json:"latitude" bson:"latitude"`
	Longitude   float64    `json:"longitude" bson:"longitude"`
	Coordinates [2]float64 `json:"coordinates" bson:"coordinates"`
}

type Species struct {
	ID          primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	SpeciesName string             `json:"species_name" bson:"species_name"`
	Category    primitive.ObjectID `json:"category_id" bson:"category"`
	Image       string             `json:"image" bson:"image"`
	Location    Location           `json:"location" bson:"location"`
}

type Category struct {
	ID           primitive.ObjectID `json:"_id" bson:"_id"`
	CategoryName string             `json:"category_name" bson:"category_name"`
}

type Animal struct {
	ID         primitive.ObjectID `json:"_id" bson:"_id"`
	AnimalName string             `json:"animal_name" bson:"animal_name"`
	Species    primitive.ObjectID `json:"species_id" bson:"species"`
	Birthdate  time.Time          `json:"birthdate" bson:"birthdate"`
	Image      string             `json:"image" bson:"image"`
	Location   Location           `json:"location" bson:"location"`
}

var (
	speciesCollection    *mongo.Collection
	categoriesCollection *mongo.Collection
	animalsCollection    *mongo.Collection
)

// @title Animal API
// @version 1.0
// @description This is a sample server for managing animals, species, and categories.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("Error loading .env file:", err)
	}

	MONGODB_URI := os.Getenv("MONGODB_URI")
	clientOptions := options.Client().ApplyURI(MONGODB_URI)
	client, err := mongo.Connect(context.Background(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")
	db := client.Database("palvelinohjelmointi")
	speciesCollection = db.Collection("species")
	categoriesCollection = db.Collection("categories")
	animalsCollection = db.Collection("animals")

	fmt.Println(speciesCollection)
	fmt.Println(categoriesCollection)
	fmt.Println(animalsCollection)

	app := fiber.New()

	// Swagger setup
	app.Get("/swagger/*", swagger.HandlerDefault)

	app.Get("/api/species", getSpecies)
	app.Get("/api/species/:id", getSpecie)
	app.Get("/api/species/name/:species_name", getSpecieByName)
	app.Post("/api/species", createSpecies)
	app.Put("/api/species/:id", updateSpecie)
	app.Delete("/api/species/:id", deleteSpecie)

	app.Get("/api/categories", getCategories)
	app.Get("/api/categories/:id", getCategory)
	app.Post("/api/categories", createCategory)
	app.Put("/api/categories/:id", updateCategory)
	app.Delete("/api/categories/:id", deleteCategory)

	app.Get("/api/animals", getAnimals)
	app.Get("/api/animals/:id", getAnimal)
	app.Post("/api/animals", createAnimal)
	app.Put("/api/animals/:id", updateAnimal)
	app.Delete("/api/animals/:id", deleteAnimal)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	log.Fatal(app.Listen("0.0.0.0:" + port))
}

// @Summary Get all species
// @Description Get a list of all species
// @Tags species
// @Produce json
// @Success 200 {array} Species
// @Router /species [get]
func getSpecies(c *fiber.Ctx) error {
	var species []Species
	ctx := context.Background()

	// Sorting
	sort := bson.D{}
	sortBy := c.Query("sort_by")
	order := -1
	if sortBy != "" {
		if c.Query("order") == "asc" {
			order = 1
		}
		sort = append(sort, bson.E{Key: sortBy, Value: order})
	}

	// Pagination
	page := c.QueryInt("page")
	limit := c.QueryInt("limit")
	skip := 0
	if page > 0 && limit > 0 {
		skip = (page - 1) * limit
	}

	findOptions := options.Find()
	if sortBy != "" {
		findOptions.SetSort(sort)
	}
	if page > 0 && limit > 0 {
		findOptions.SetSkip(int64(skip))
		findOptions.SetLimit(int64(limit))
	}

	cursor, err := speciesCollection.Find(ctx, bson.M{}, findOptions)
	if err != nil {
		fmt.Printf("Error fetching species: %v\n", err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch species"})
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var specie Species
		if err := cursor.Decode(&specie); err != nil {
			fmt.Printf("Error decoding species: %v\n", err)
			return c.Status(500).JSON(fiber.Map{"error": "Failed to decode species"})
		}

		species = append(species, specie)
	}

	if err := cursor.Err(); err != nil {
		fmt.Printf("Cursor error: %v\n", err)
		return c.Status(500).JSON(fiber.Map{"error": "Cursor error"})
	}

	return c.JSON(species)
}

// @Summary Create a new species
// @Description Create a new species
// @Tags species
// @Accept json
// @Produce json
// @Param species body Species true "Species to create"
// @Success 201 {object} Species
// @Router /species [post]
func createSpecies(c *fiber.Ctx) error {
	species := new(Species)

	if err := c.BodyParser(species); err != nil {
		return err
	}
	if species.SpeciesName == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Species name is required"})
	}
	if species.Category.IsZero() {
		return c.Status(400).JSON(fiber.Map{"error": "Category ID is required"})
	}

	species.ID = primitive.NewObjectID()

	_, err := speciesCollection.InsertOne(context.Background(), species)
	if err != nil {
		return err
	}
	return c.Status(201).JSON(species)
}

// @Summary Update a species
// @Description Update a species by its ID
// @Tags species
// @Accept json
// @Produce json
// @Param id path string true "Species ID"
// @Param species body Species true "Species to update"
// @Success 200 {object} Species
// @Router /species/{id} [put]
func updateSpecie(c *fiber.Ctx) error {
	id := c.Params("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}

	var updateData map[string]interface{}
	if err := c.BodyParser(&updateData); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}
	update := bson.M{
		"$set": updateData,
	}

	filter := bson.M{"_id": objID}
	ctx := context.Background()

	result, err := speciesCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to update species"})
	}

	return c.JSON(result)
}

// @Summary Get a species by ID
// @Description Get a species by its ID
// @Tags species
// @Produce json
// @Param id path string true "Species ID"
// @Success 200 {object} Species
// @Router /species/{id} [get]
func getSpecie(c *fiber.Ctx) error {
	id := c.Params("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}

	var specie Species
	ctx := context.Background()

	err = speciesCollection.FindOne(ctx, bson.M{"_id": objID}).Decode(&specie)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(404).JSON(fiber.Map{"error": "Species not found"})
		}
		return c.Status(500).JSON(fiber.Map{"error": "Failed to retrieve species"})
	}

	return c.JSON(specie)
}

// @Summary Get a species by name
// @Description Get a species by its name
// @Tags species
// @Produce json
// @Param species_name path string true "Species Name"
// @Success 200 {object} Species
// @Router /species/name/{species_name} [get]
func getSpecieByName(c *fiber.Ctx) error {
	var specie Species
	ctx := context.Background()

	speciesName := c.Params("species_name")

	filter := bson.M{"species_name": speciesName}
	err := speciesCollection.FindOne(ctx, filter).Decode(&specie)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(404).JSON(fiber.Map{"error": "Species not found"})
		}
		fmt.Printf("Error fetching species: %v\n", err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch species"})
	}

	return c.JSON(specie)
}

// @Summary Delete a species
// @Description Delete a species by its ID
// @Tags species
// @Param id path string true "Species ID"
// @Success 204
// @Router /species/{id} [delete]
func deleteSpecie(c *fiber.Ctx) error {
	id := c.Params("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}

	filter := bson.M{"_id": objID}
	ctx := context.Background()

	_, err = speciesCollection.DeleteOne(ctx, filter)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to delete species"})
	}
	return c.Status(200).JSON(fiber.Map{"success": true})
}

// @Summary Get a category by ID
// @Description Get a category by its ID
// @Tags categories
// @Produce json
// @Param id path string true "Category ID"
// @Success 200 {object} Category
// @Router /categories/{id} [get]
func getCategory(c *fiber.Ctx) error {
	id := c.Params("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}

	var category Category
	ctx := context.Background()

	err = categoriesCollection.FindOne(ctx, bson.M{"_id": objID}).Decode(&category)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(404).JSON(fiber.Map{"error": "Category not found"})
		}
		return c.Status(500).JSON(fiber.Map{"error": "Failed to retrieve category"})
	}

	return c.JSON(category)
}

// Categories routes
// @Summary Get all categories
// @Description Get a list of all categories
// @Tags categories
// @Produce json
// @Success 200 {array} Category
// @Router /categories [get]
func getCategories(c *fiber.Ctx) error {
	var categories []Category
	ctx := context.Background()

	// Sorting
	sort := bson.D{}
	sortBy := c.Query("sort_by")
	order := -1
	if sortBy != "" {
		if c.Query("order") == "asc" {
			order = 1
		}
		sort = append(sort, bson.E{Key: sortBy, Value: order})
	}

	// Pagination
	page := c.QueryInt("page")
	limit := c.QueryInt("limit")
	skip := 0
	if page > 0 && limit > 0 {
		skip = (page - 1) * limit
	}

	findOptions := options.Find()
	if sortBy != "" {
		findOptions.SetSort(sort)
	}
	if page > 0 && limit > 0 {
		findOptions.SetSkip(int64(skip))
		findOptions.SetLimit(int64(limit))
	}

	cursor, err := categoriesCollection.Find(ctx, bson.M{}, findOptions)
	if err != nil {
		fmt.Printf("Error fetching categories: %v\n", err) // Debugging statement
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch categories"})
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var category Category
		if err := cursor.Decode(&category); err != nil {
			fmt.Printf("Error decoding category: %v\n", err) // Debugging statement
			return c.Status(500).JSON(fiber.Map{"error": "Failed to decode category"})
		}
		categories = append(categories, category)
	}

	if err := cursor.Err(); err != nil {
		fmt.Printf("Cursor error: %v\n", err) // Debugging statement
		return c.Status(500).JSON(fiber.Map{"error": "Cursor error"})
	}

	return c.JSON(categories)
}

// @Summary Create a new category
// @Description Create a new category
// @Tags categories
// @Accept json
// @Produce json
// @Param category body Category true "Category to create"
// @Success 201 {object} Category
// @Router /categories [post]
func createCategory(c *fiber.Ctx) error {
	category := new(Category)

	if err := c.BodyParser(category); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}
	if category.CategoryName == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Category name is required"})
	}

	category.ID = primitive.NewObjectID()

	_, err := categoriesCollection.InsertOne(context.Background(), category)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create category"})
	}
	return c.Status(201).JSON(category)
}

// @Summary Update a category
// @Description Update a category by its ID
// @Tags categories
// @Accept json
// @Produce json
// @Param id path string true "Category ID"
// @Param category body Category true "Category to update"
// @Success 200 {object} Category
// @Router /categories/{id} [put]
func updateCategory(c *fiber.Ctx) error {
	id := c.Params("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}

	var updateData map[string]interface{}
	if err := c.BodyParser(&updateData); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}
	update := bson.M{
		"$set": updateData,
	}

	filter := bson.M{"_id": objID}

	result, err := categoriesCollection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to update document"})
	}

	return c.JSON(result)
}

// @Summary Delete a category
// @Description Delete a category by its ID
// @Tags categories
// @Param id path string true "Category ID"
// @Success 204
// @Router /categories/{id} [delete]
func deleteCategory(c *fiber.Ctx) error {
	id := c.Params("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}

	filter := bson.M{"_id": objID}
	_, err = categoriesCollection.DeleteOne(context.Background(), filter)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to delete category"})
	}
	return c.Status(200).JSON(fiber.Map{"success": true})
}

// Animals routes
// @Summary Get all animals
// @Description Get a list of all animals
// @Tags animals
// @Produce json
// @Success 200 {array} Animal
// @Router /animals [get]
func getAnimals(c *fiber.Ctx) error {
	ctx := context.Background()

	// Sorting
	sort := bson.D{}
	sortBy := c.Query("sort_by")
	order := 1
	if sortBy != "" {
		if c.Query("order") == "desc" {
			order = -1
		}
		sort = append(sort, bson.E{Key: sortBy, Value: order})
	}

	// Pagination
	page := c.QueryInt("page")
	limit := c.QueryInt("limit")
	skip := 0
	if page > 0 && limit > 0 {
		skip = (page - 1) * limit
	}

	// Define the aggregation pipeline
	pipeline := mongo.Pipeline{
		{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "species"},
			{Key: "localField", Value: "species"},
			{Key: "foreignField", Value: "_id"},
			{Key: "as", Value: "species"},
		}}},
		{{Key: "$unwind", Value: bson.D{
			{Key: "path", Value: "$species"},
			{Key: "preserveNullAndEmptyArrays", Value: true},
		}}},
		{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "categories"},
			{Key: "localField", Value: "species.category"},
			{Key: "foreignField", Value: "_id"},
			{Key: "as", Value: "species.category"},
		}}},
		{{Key: "$unwind", Value: bson.D{
			{Key: "path", Value: "$species.category"},
			{Key: "preserveNullAndEmptyArrays", Value: true},
		}}},
	}

	// Add sorting stage if sort key is present
	if len(sort) > 0 {
		pipeline = append(pipeline, bson.D{{Key: "$sort", Value: sort}})
	}

	// Add pagination stages if page and limit are present
	if page > 0 && limit > 0 {
		pipeline = append(pipeline, bson.D{{Key: "$skip", Value: skip}})
		pipeline = append(pipeline, bson.D{{Key: "$limit", Value: int64(limit)}})
	}

	// Execute the aggregation pipeline
	cursor, err := animalsCollection.Aggregate(ctx, pipeline)
	if err != nil {
		fmt.Printf("Error executing aggregation pipeline: %v\n", err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch animals"})
	}
	defer cursor.Close(ctx)

	var animals []bson.M
	if err := cursor.All(ctx, &animals); err != nil {
		fmt.Printf("Error decoding animals: %v\n", err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to decode animals"})
	}

	return c.JSON(animals)
}

// @Summary Get an animal by ID
// @Description Get an animal by its ID
// @Tags animals
// @Produce json
// @Param id path string true "Animal ID"
// @Success 200 {object} Animal
// @Router /animals/{id} [get]
func getAnimal(c *fiber.Ctx) error {
	id := c.Params("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}

	var animal Animal
	err = animalsCollection.FindOne(context.Background(), bson.M{"_id": objID}).Decode(&animal)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(404).JSON(fiber.Map{"error": "Animal not found"})
		}
		return c.Status(500).JSON(fiber.Map{"error": "Failed to retrieve animal"})
	}

	return c.JSON(animal)
}

// @Summary Create a new animal
// @Description Create a new animal
// @Tags animals
// @Accept json
// @Produce json
// @Param animal body Animal true "Animal to create"
// @Success 201 {object} Animal
// @Router /animals [post]
func createAnimal(c *fiber.Ctx) error {
	animal := new(Animal)

	if err := c.BodyParser(animal); err != nil {
		return err
	}
	if animal.AnimalName == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Animal name is required"})
	}

	animal.ID = primitive.NewObjectID()

	_, err := animalsCollection.InsertOne(context.Background(), animal)
	if err != nil {
		return err
	}
	return c.Status(201).JSON(animal)
}

// @Summary Update an animal
// @Description Update an animal by its ID
// @Tags animals
// @Accept json
// @Produce json
// @Param id path string true "Animal ID"
// @Param animal body Animal true "Animal to update"
// @Success 200 {object} Animal
// @Router /animals/{id} [put]
func updateAnimal(c *fiber.Ctx) error {
	id := c.Params("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}

	var updateData map[string]interface{}
	if err := c.BodyParser(&updateData); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}
	update := bson.M{
		"$set": updateData,
	}

	filter := bson.M{"_id": objID}

	result, err := animalsCollection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to update document"})
	}

	return c.JSON(result)
}

// @Summary Delete an animal
// @Description Delete an animal by its ID
// @Tags animals
// @Param id path string true "Animal ID"
// @Success 204
// @Router /animals/{id} [delete]
func deleteAnimal(c *fiber.Ctx) error {
	id := c.Params("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}

	filter := bson.M{"_id": objID}
	_, err = animalsCollection.DeleteOne(context.Background(), filter)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to delete animal"})
	}
	return c.Status(200).JSON(fiber.Map{"success": true})
}
