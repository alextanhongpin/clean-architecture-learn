# Constructor 


Use constructor for creating valid only object. If they cannot be created as valid only object (may because type can be initialize without calling constructor), then they should have a validate method instead.
https://stackoverflow.com/questions/43456801/best-practices-for-long-constructors-in-javascript


If you have validation logic that needs to be shared between creation and update, then place the logic in a separate function or method of the class.
