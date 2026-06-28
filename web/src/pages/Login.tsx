import { supabase } from "../lib/supabase";
import { useNavigate } from "react-router-dom";

function Login() {

const navigate = useNavigate();

const handleLogin = async () => {
  const { data, error } = await supabase.auth.signInWithPassword({
    email: "rbac-test@test.com",
    password: "password123",
  });
  if (error) {
    console.error("login error:", error.message);
    return;
  }
  console.log("logged in:", data.session?.user.id);
  navigate("/dashboard");
};

  return (
    <div>
      <h1>IssueBot</h1>
      <p>Login to continue</p>
      <button onClick={handleLogin}>Login</button>
    </div>
  );
}

export default Login;