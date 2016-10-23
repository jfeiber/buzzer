package main

import (
  "math/rand"
  "fmt"
)

var adjectives = []string{"adorable", "beautiful", "clean", "drab", "elegant", "fancy", "plain",
                          "clever", "gifted", "helpful", "odd", "shy", "brave", "calm", "delightful",
                          "eager", "goofy", "sneaky", "gentle", "happy", "jolly", "kind", "proud",
                          "silly", "witty", "angry", "clumsy", "turnt", "fierce", "grumpy", "lazy",
                          "mysterious", "nervous", "obnoxious", "chubby", "skinny", "ancient", "bitter",
                          "strong", "cold", "cool", "wet"}

var animals = []string{"aardvark", "albatross", "alligator", "alpaca", "ant", "anaconda", "anteater",
                       "antelope", "ape", "baboon", "badger", "barracuda", "bear", "beaver", "eagle",
                       "bee", "bird", "bison", "bobcat", "bonobo", "buffalo", "butterfly", "camel",
                       "cat", "caterpillar", "chameleon", "cheetah", "chicken", "chimapnzee",
                       "chimpmunk", "cobra", "cougar", "cow", "coyote", "crab", "crane", "cricket",
                       "crocodile", "cricket", "crow", "deer", "dinosaur", "dog", "dolphin", "donkey",
                       "dove", "dragonfly", "dragon", "duck", "eel", "elephant", "elk", "falcon",
                       "firefly", "fish", "flamingo", "fox", "frog", "gazelle", "gecko", "panda",
                       "squid", "giraffe", "goat", "goldfish", "goose", "gopher", "gorilla",
                       "grasshopper", "heron", "shark", "sloth", "pig", "hamster", "hawk", "hornet",
                       "horse", "hummingbird", "iguana", "impala", "jaguar", "jellyfish", "kangaroo",
                       "lemur", "lion", "leopard", "lizard", "llama", "lobster", "manatee",
                       "marlin", "meerkat", "monkey", "mosquito", "moth", "mouse", "narwhal",
                       "newt", "ocelot", "octopus", "orangutan", "otter", "ostrich", "owl", "ox",
                       "panther", "parrot", "pelican", "penguin", "pheasant", "pirahana", "porcupine",
                       "puma", "python", "quail", "rabbit", "raven", "snake", "reindeer", "rhino",
                       "rooster", "salamander", "salmon", "scorpion", "seahorse", "sparrow", "spider",
                       "squirrel", "swan", "swordfish", "tiger", "toad", "tortoise", "tuna", "turkey",
                       "turtle", "walrus", "whale", "wolf", "worm", "yak", "zebra", "kevin"}

type BuzzerNameGenerator interface {
  GenerateName() string
}

type RandomBuzzerName struct {
  r *rand.Rand
}

func (randBuzzerName RandomBuzzerName) GenerateName() string {
  r := randBuzzerName.r
  return fmt.Sprintf("%s-%s-%d", adjectives[r.Intn(len(adjectives))], animals[r.Intn(len(animals))],
                     r.Intn(10000))
}

func NewBuzzerNameGenerator(seed int64) BuzzerNameGenerator {
  generator := RandomBuzzerName{rand.New(rand.NewSource(seed))}
  return generator
}
