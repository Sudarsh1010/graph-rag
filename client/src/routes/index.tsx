import { createFileRoute } from "@tanstack/react-router"

// eslint-disable-next-line react-refresh/only-export-components
function Index() {
  return (
    <div className="p-2">
      <h3>Welcome Home!</h3>
    </div>
  )
}

export const Route = createFileRoute("/")({
  component: Index,
})
