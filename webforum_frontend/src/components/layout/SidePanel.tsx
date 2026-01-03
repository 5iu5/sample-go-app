import {useState, useEffect} from 'react'
import type {Topic} from "../../Api"
import "./Layout.css"
export default function SidePanel(){
    const [topics, setTopics] = useState<Topic[] | null> (null)
    
    useEffect(()=>{
        async function loadTopics(){
            try{
                const res = await fetch("http://localhost:8000/topics", {
                    
                    credentials: 'include',
                    
                })
                const data = await res.json()
                setTopics(data)

            } catch(err){
                console.log("Failed to load topics", err)
            }
        }
        loadTopics()
        console.log("Sucessfully called loadTopic")
    }, [])
    return(
        <div style={{display:"flex", flexDirection: "row"}}>
            <div style={{display:"flex", flexDirection: "column",minWidth:"13vw", maxWidth:"20vw", paddingTop:"10px", paddingBottom:"10px"}}>
                

                <div style={{padding:"0px", marginTop:"20%"}}>
                    <p style={{margin:0, marginLeft:"10px", padding:0, color:"#818080ff", fontSize:"1rem"}}>Topics</p>
                    <div style={{width:"95%", justifySelf:"center", height:"1px", backgroundColor:"rgba(153, 151, 150, 1)"}}/>
                    <ul style={{ listStyleType: "none", margin:"0", padding:"0"}}>
                    {topics && topics.map(topic => (  
                        <li key={topic.topic_id} className="topic-item"
                        onClick={()=>{console.log("topic name is: ", topic.name)}}>
                            {topic.name}
                        </li>                   
                    ))}
                    </ul>
                </div>

                <div style={{width:"90%", alignSelf:"center", height:"1px", backgroundColor:"rgba(153, 151, 150, 1)"}}/>


                
            </div>
        
            <div style={{ width:"1px", height:"100vh", backgroundColor:"#1d1d1dff"}}/>
        </div>
    )


}