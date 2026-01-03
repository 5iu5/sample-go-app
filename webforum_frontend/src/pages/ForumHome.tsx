import PostCard from "../components/PostCard"
import TopicCard from "../components/TopicCard"
import {useState, useEffect} from "react"
import type {Post} from "../Api"
import {getAllPosts} from "../Api"

export default function ForumHome() {
    //temporary
    const topics = ["Tech", "Games", "LifeStyle", "Music", "Automotive", "Culture", "Sports", "Farming", "Programming", "Science"]
    const [posts, setPosts] = useState<Post[] | null>(null)

    useEffect(() => {
        console.log("running useEffect")
        async function loadPosts(){
            try {
                const data = await getAllPosts()
                setPosts(data)
            } catch (err){
                console.log("failed retrieving posts", err)
            }
        }
        loadPosts();
    }, []);

    return(
        <div style={{ display: "flex", flexDirection: "column", width:"100%", alignItems:"center", gap:"15px", height:"100vh", overflowY: "auto", marginBottom: "10px"}}>
            
            {/* <h1 style={{alignSelf: "center"}}>Forum Topics </h1>

             {topics.map((topic) => (
                <TopicCard key={topic} title = {topic}/>
             ))} */}

            {posts && posts.map(post => (
                <PostCard key={post.post_id} post={post}/>
            ))}
        </div>
    )
}