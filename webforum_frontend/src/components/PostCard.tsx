import '../css/Forum.css'
import type { Post } from "../Api";
import {getPostCommentCount} from "../Api"
import {useState, useEffect} from "react"
import {Link, useNavigate} from "react-router-dom"
import CommentIcon from '../assets/icons/CommentIcon'
import UpvoteIcon from '../assets/icons/UpvoteIcon'
import DownvoteIcon from '../assets/icons/DownvoteIcon'
import calcDuration from '../functions/calcDuration'

export default function PostCard({post} : PostCardProps) {
    const navigate = useNavigate()

    
    const [commentCount, setCommentCount] = useState(0)
    useEffect(()=>{
        async function loadCommentCount(){
            try{
                const data = await getPostCommentCount(post.post_id)
                setCommentCount(data)
            }
            catch(err){
                console.log("Error getting no. of comments for post", err)
            }
        }
        loadCommentCount()
    },)


    return(
        
        <div className = "postCardDiv"
         onClick={() => {
            navigate(`/postpage/${post.post_id}`)
            console.log("Post clicked!");
        }}>
            
            <div style = {{display: "flex",flexDirection: "row", backgroundColor:"lightgreen", width:"fit-content"}}
                onClick={(e)=>{
                    e.stopPropagation();        //stops propagation to prevent outer div's onClick from running
                    console.log("username/profile clicked!")
                }}> 
                <Link to={`/user`} style={{color:"black",display:"flex", flexDirection:"row"}}>
                <span style ={{fontSize: "13px"}}> profile picture & username</span>
                </Link>
                <div style={{marginRight: "10px"}}/>
                <span style={{ fontSize: "12px", color: "gray" }}>{calcDuration(post.created_at)}</span>
                
            </div>    
   
            <div className="postContent"
                onClick={(e) => {
                console.log("content clicked!");
            }}>
                <h3 className="titleText">{post.text_title}</h3>    
                <p className="bodyText">{post.text_body}</p>
            </div>



            <div style={{display:"flex", flexDirection:"row", alignSelf: "flex-start", alignItems:"center", fontSize: "12px",  color: "gray", height:"20px", gap:"8px",  marginTop:"2px", marginLeft:"5px" }}
                onClick={(e) => {
                    e.stopPropagation(); 
                    console.log("comments/upvotes clicked!");
                }}>
                
                <div style={{display:"flex", flexDirection:"row",alignItems:"center", gap:"2px"}}>
                <p style={{fontSize: "13px",padding:"0", margin:"0"}}>0</p>
                <UpvoteIcon/>
                </div>
                
                <div style={{display:"flex", flexDirection:"row",alignItems:"center", gap:"2px"}}>
                <p style={{fontSize: "13px",padding:"0", margin:"0"}}>0</p>
                <DownvoteIcon/>
                </div>

                <div style={{display:"flex", flexDirection:"row", alignItems:"center", gap:"2px"}}>
                <p style={{fontSize: "13px",padding:"0", margin:"0"}}>{commentCount}</p>
                <CommentIcon/>
                </div>
                
                </div>
        </div>

        
        
    )
}

interface PostCardProps {
  post: Post;
}