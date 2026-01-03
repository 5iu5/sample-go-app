import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import './index.css'
import Layout from './components/layout/Layout.tsx'
import PostPage from './pages/PostPage.tsx'
import ForumHome from './pages/ForumHome.tsx'
import LoginPage from './pages/LoginPage.js'
import { createBrowserRouter, RouterProvider } from 'react-router-dom'

const router = createBrowserRouter([
  {
    path:"/",
    element: <Layout/>,
    children: [
      {path:"/", element: <ForumHome/>},
      {path:"/postpage/:id", element: <PostPage/>},
      // {path:"/topic/:id", element: <TopicPage/>}

      // {path:"*", element:<NotFoundPage/>}
    ]
  },
  {
    path: "/login",
    element: <LoginPage/>
  }
])

createRoot(document.getElementById('root')!).render(
  <StrictMode>
    <RouterProvider router = {router}/>
    {/* <App /> */}
  </StrictMode>,
)
