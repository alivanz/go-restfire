package restfire

func NewApp(key string) FirebaseAuthProvider {
	return &firebaseAuth{key}
}
