namespace go thrift.rpc

struct User {
    1: required i32 uid;                
    2: required string name;
    5: optional Profile pro;
}

struct Profile {                  
	1: required i32 uid;
    2: required i16 age;
}

service userService {        
		i64 AddUser(1:string name, 2:i16 age),
		User GetUser(1:i32 uid),
		list<User> GetAllUsers(1:i32 rows, 2:i32 page),
		User UpdateUser(1:i32 uid, 2:string name, 3:i16 age),
		i64 DeleteUser(1:i32 uid),
}