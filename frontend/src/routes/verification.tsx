import { createFileRoute } from '@tanstack/react-router'

export const Route = createFileRoute('/verification')({
  component: RouteComponent,
})

function RouteComponent() {
  return <div>Hello "/verification"!</div>
}
