import { createContext, useContext, useState } from 'react'

const AuthContext = createContext(null)

const AuthProvider = ({children}) => {
  const [signInState, setSignInState] = useState(false)

  const isSignIn = () => {
    return signInState
  }

  const authSignIn = () => {
    setSignInState(true)
  };

  const authSignOut = () => {
    setSignInState(false)
  };

  const auth = {
    isSignIn,
    authSignIn,
    authSignOut,
  };
  return (
    <AuthContext.Provider value={auth}>
      {children}
    </AuthContext.Provider>
  );
}

function useAuthContext() {
  const auth = useContext(AuthContext)
  if (auth === null) {
    throw new Error(
      'Component must be wrapped in Provider in order to access API',
    );
  }
  return auth
}

export { AuthProvider, useAuthContext }