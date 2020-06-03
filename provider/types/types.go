package types

//Record is the exchange type between Loader and Saver.
type Record map[string]string

//NullValue flag string for NULL value in database columns.
const NullValue string = "NULL&NULL@NIL"
