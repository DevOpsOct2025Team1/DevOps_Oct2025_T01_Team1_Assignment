import { useEffect } from "react"
import { useNavigate } from "react-router"
import { useAuth } from "~/contexts/AuthContext"

const Logout = () => {
  const navigate = useNavigate()
  const { clearAuth } = useAuth()

  useEffect(() => {
    clearAuth()
    navigate("/login")
  }, [clearAuth, navigate])

  return null
}

export default Logout
