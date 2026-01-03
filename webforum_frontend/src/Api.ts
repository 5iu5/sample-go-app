const API = "http://localhost:8000"

/*Users*/

// List all users
export async function fetchUsers() {
  const res = await fetch(`${API}/users`, {
    method: "GET",
    credentials: "include"
  });
  return res.json();
}

// Get a user by ID
export async function fetchUser(userId: number) {
  const res = await fetch(`${API}/users/${userId}`, {
    method: "GET",
    credentials: "include"
  });
  return res.json();
}

// Get posts by a user
export async function fetchUserPosts(userId: number) {
  const res = await fetch(`${API}/users/${userId}/posts`, {
    method: "GET",
    credentials: "include"
  });
  return res.json();
}

// Get comments by a user
export async function fetchUserComments(userId: number) {
  const res = await fetch(`${API}/users/${userId}/comments`, {
    method: "GET",
    credentials: "include"
  });
  return res.json();
}


/* Post */
//get all posts
export async function getAllPosts(){
  console.log("running get post funct")
  const res = await fetch(`${API}/posts`, {
    method: "GET",
    credentials: "include"
  })
  return res.json()
}

export async function getPostById(id: number){     //id is post id
  const res = await fetch(`${API}/posts/${id}`, {
    method: "GET",
    credentials: "include"
  })
  return res.json()       //.json() converts the json array to javascript array
}

export async function addPost(data: AddPostForm){
  const res = await fetch(`${API}/posts`, {
      method: "POST",
      headers: {"Content-Type": "application/json"},
      body: JSON.stringify(data),
      credentials: "include"
  })
  return res.json();
}

export async function editPost(id: number, data: EditPostForm){
  const res = await fetch(`${API}/posts/${id}`, {
      method: "PATCH",
      headers: {"Content-Type": "application/json"},
      body: JSON.stringify(data),
      credentials: "include"
  })
  return res.json()

}

export async function deletePost(id: number){
  const res = await fetch(`${API}/posts/${id}`, {
    method: "DELETE",
    credentials: "include"
  })
  return res.json()
}

/* Comments */
export async function getPostComments(id: number){
  const res = await fetch(`${API}/posts/${id}/comments`, {
    method: "GET",
    credentials: "include"
  })
  return res.json()
}

export async function getPostCommentCount(id: number){
  const res = await fetch(`${API}/posts/${id}/comments/count`, {
    method: "GET",
    credentials: "include"
  })
  return res.json()
}

export async function addComment(id: number, data: AddCommentForm){
    const res = await fetch(`${API}/posts/${id}/comments`, {
        method: "POST",
        headers: {"Content-Type": "application/json"},
        body: JSON.stringify(data),
        credentials: "include"
    })
    return res.json()
}

export async function editComment(id: number, data: EditCommentForm){
    const res = await fetch(`${API}/comments/${id}`, {
        method: "PATCH",
        headers: {"Content-Type": "application/json"},
        body: JSON.stringify(data),
        credentials: "include"
    })
    return res.json()
}
export async function deleteComment(id: number){
    const res = await fetch(`${API}/comments/${id}`, { 
        method: "DELETE",
        credentials: "include"
    });
    return res.json()
}

export async function userLogin(data: loginForm){
  const res = await fetch(`${API}/auth/login`, {
    method:"POST",
    headers: {"Content-Type": "application/json"},
    body: JSON.stringify(data),
    credentials: "include"
  })
  return res.json()
}

export type loginForm = {
  username: string
  password: string
}

interface AddPostForm {
  topic_id: number
  user_id: number
  text_title: string
  text_body: string
}
interface EditPostForm {
  text_body: string
}

interface AddCommentForm {
  post_id: number
  parent_id?: number | null
  user_id: number
  text_body: string
}
interface EditCommentForm {
  text_body: string
}

// interface AddUserForm {
//   user_id: number;
//   username: string;
//   email: string;
// }

export interface Post {
  post_id: number
  topic_id: number
  user_id: number
  text_title: string
  text_body: string
  created_at: string
  is_deleted: boolean
}

export interface Comment {
  comment_id: number
  post_id: number
  parent_id: number
  user_id: number
  text_body: string 
  created_at: string 
  updated_at?: string
  is_deleted: boolean
  username: string
}

export interface Topic{
  topic_id: number
  name: string
  description: string | null
  created_at: string
}