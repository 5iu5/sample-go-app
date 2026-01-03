import { userLogin } from "../Api"
import { useNavigate } from "react-router-dom"
import "./LoginPage.css"
export default function LoginPage(){
    const navigate = useNavigate()
    async function submitHandler(event: React.FormEvent<HTMLFormElement>){
        event.preventDefault()
        const form = event.currentTarget
        const formData = new FormData(form)

        const username = formData.get("username") as string
        const password = formData.get("password") as string
        

        try{
            const response: responseType = await userLogin({username: username, password: password})
            console.log(response)
            if (response.success == true){
                navigate("/")
            }
            else{
                alert("Incorrect username/password")
            }
        } catch(err){
            console.log("login failed: ", err)
        }

        
    }

    return(
        <>
            <div className="formBox">
                <h3>Welcome to Forum's Login Page</h3>
                <form  onSubmit={(event) =>{
                    submitHandler(event)
                }}>
                    <div>
                        <label>username </label>
                        <input name="username" type="text" required/>
                    </div>

                    <div>
                        <label> password </label>
                        <input name="password" type="password"/>
                        <input type="submit"/>
                    </div>

                    

                </form>


            </div>
            </>
                

        
    )
}

type responseType = {
    success: boolean, 
    userId: number, 
    message: string
}