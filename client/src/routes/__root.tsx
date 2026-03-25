import { Link, Outlet, createRootRoute, useMatchRoute } from "@tanstack/react-router"
import {
  HouseIcon,
  DatabaseIcon,
  GraphIcon,
  MagnifyingGlassIcon,
  ChatCircleIcon,
} from "@phosphor-icons/react"
import { TooltipProvider } from "~/components/ui/tooltip"
import {
  Sidebar,
  SidebarContent,
  SidebarFooter,
  SidebarGroup,
  SidebarGroupContent,
  SidebarGroupLabel,
  SidebarHeader,
  SidebarInset,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
  SidebarProvider,
  SidebarRail,
  SidebarTrigger,
} from "~/components/ui/sidebar"

const navItems = [
  { to: "/", label: "Home", icon: HouseIcon },
  { to: "/data", label: "Data Explorer", icon: DatabaseIcon },
  { to: "/graph", label: "Graph", icon: GraphIcon },
  { to: "/search", label: "Search", icon: MagnifyingGlassIcon },
  { to: "/ask", label: "Ask", icon: ChatCircleIcon },
]

export const Route = createRootRoute({
  component: () => {
    const matchRoute = useMatchRoute()
    
    return (
    <TooltipProvider>
      <SidebarProvider>
        <Sidebar collapsible="icon">
          <SidebarHeader className="border-b border-sidebar-border">
            <SidebarMenu>
              <SidebarMenuItem>
                <SidebarMenuButton
                  size="lg"
                  isActive={!!matchRoute({ to: "/" })}
                  render={<Link to="/" />}
                >
                  <div className="flex aspect-square size-8 items-center justify-center rounded-lg bg-primary text-primary-foreground">
                    <GraphIcon className="size-4" />
                  </div>
                  <div className="grid flex-1 text-left text-sm leading-tight">
                    <span className="truncate font-semibold">Graph RAG</span>
                  </div>
                </SidebarMenuButton>
              </SidebarMenuItem>
            </SidebarMenu>
          </SidebarHeader>
          <SidebarContent>
            <SidebarGroup>
              <SidebarGroupLabel>Navigation</SidebarGroupLabel>
              <SidebarGroupContent>
                <SidebarMenu>
                  {navItems.map((item) => {
                    const isActive = !!matchRoute({ to: item.to })
                    
                    return (
                      <SidebarMenuItem key={item.to}>
                        <SidebarMenuButton
                          tooltip={item.label}
                          isActive={isActive}
                          render={<Link to={item.to} />}
                        >
                          <item.icon />
                          <span>{item.label}</span>
                        </SidebarMenuButton>
                      </SidebarMenuItem>
                    )
                  })}
                </SidebarMenu>
              </SidebarGroupContent>
            </SidebarGroup>
          </SidebarContent>
          <SidebarFooter className="border-t border-sidebar-border">
            <SidebarMenu>
              <SidebarMenuItem>
                <SidebarMenuButton size="sm" className="text-xs text-sidebar-foreground/70">
                  <span className="truncate">Graph RAG Client</span>
                </SidebarMenuButton>
              </SidebarMenuItem>
            </SidebarMenu>
          </SidebarFooter>
          <SidebarRail />
        </Sidebar>
        <SidebarInset>
          <header className="flex h-12 shrink-0 items-center gap-2 border-b px-4">
            <SidebarTrigger />
          </header>
          <main className="flex-1 overflow-auto p-4">
            <Outlet />
          </main>
        </SidebarInset>
      </SidebarProvider>
    </TooltipProvider>
  )
  },
})
