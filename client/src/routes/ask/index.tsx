import { createFileRoute } from "@tanstack/react-router"
import { useState, useRef, useEffect } from "react"
import { useMutation } from "@tanstack/react-query"
import { ChatCircleIcon, PaperPlaneTiltIcon, SpinnerGapIcon } from "@phosphor-icons/react"
import { Card, CardHeader, CardTitle, CardDescription, CardContent } from "~/components/ui/card"
import { Button } from "~/components/ui/button"
import { Textarea } from "~/components/ui/textarea"
import { ScrollArea } from "~/components/ui/scroll-area"
import { apiFetch } from "~/lib/api-client"

interface NLQueryResponse {
  answer: string
  data?: Record<string, unknown>[]
  visualization?: Record<string, unknown>
  sources?: string[]
}

interface Message {
  id: string
  role: "user" | "assistant"
  content: string
  data?: Record<string, unknown>[]
  visualization?: Record<string, unknown>
  sources?: string[]
}

function Ask() {
  const [input, setInput] = useState("")
  const [messages, setMessages] = useState<Message[]>([])
  const scrollRef = useRef<HTMLDivElement>(null)
  const textareaRef = useRef<HTMLTextAreaElement>(null)

  const mutation = useMutation({
    mutationFn: (question: string) =>
      apiFetch<NLQueryResponse>("/api/v1/nlquery", {
        method: "POST",
        body: JSON.stringify({ question }),
      }),
    onSuccess: (data, _variables) => {
      const assistantMessage: Message = {
        id: crypto.randomUUID(),
        role: "assistant",
        content: data.answer,
        data: data.data,
        visualization: data.visualization,
        sources: data.sources,
      }
      setMessages((prev) => [...prev, assistantMessage])
    },
    onError: (error, _variables) => {
      const errorMessage: Message = {
        id: crypto.randomUUID(),
        role: "assistant",
        content: `Error: ${error instanceof Error ? error.message : "Failed to get response"}`,
      }
      setMessages((prev) => [...prev, errorMessage])
    },
  })

  const handleSubmit = (e?: React.FormEvent) => {
    e?.preventDefault()
    if (!input.trim() || mutation.isPending) return

    const userMessage: Message = {
      id: crypto.randomUUID(),
      role: "user",
      content: input.trim(),
    }
    setMessages((prev) => [...prev, userMessage])
    mutation.mutate(input.trim())
    setInput("")
  }

  const handleKeyDown = (e: React.KeyboardEvent<HTMLTextAreaElement>) => {
    if (e.key === "Enter" && !e.shiftKey) {
      e.preventDefault()
      handleSubmit()
    }
  }

  // Auto-scroll to bottom when messages change
  useEffect(() => {
    if (scrollRef.current) {
      scrollRef.current.scrollTo({
        top: scrollRef.current.scrollHeight,
        behavior: "smooth",
      })
    }
  }, [messages])

  // Focus textarea on mount
  useEffect(() => {
    textareaRef.current?.focus()
  }, [])

  return (
    <div className="flex h-full flex-col gap-6">
      <div className="flex flex-col gap-2">
        <h1 className="font-heading text-2xl font-bold tracking-tight">Ask</h1>
        <p className="text-muted-foreground">
          Ask natural language questions about your SAP data.
        </p>
      </div>

      {/* Chat area */}
      <Card className="flex flex-1 flex-col overflow-hidden">
        <CardHeader className="border-b py-3">
          <CardTitle className="flex items-center gap-2 text-base">
            <ChatCircleIcon className="size-5 text-primary" />
            Chat
          </CardTitle>
          <CardDescription>
            Type your question below. Press Enter to send, Shift+Enter for newline.
          </CardDescription>
        </CardHeader>
        <CardContent className="flex flex-1 flex-col p-0">
          {/* Messages */}
          <ScrollArea className="flex-1 p-4" ref={scrollRef}>
            <div className="flex flex-col gap-4">
              {messages.length === 0 && (
                <div className="flex flex-col items-center justify-center py-12 text-center">
                  <ChatCircleIcon className="size-12 text-muted-foreground" />
                  <div className="mt-4 text-lg font-medium">Start a conversation</div>
                  <div className="mt-1 text-sm text-muted-foreground">
                    Ask questions like "Show me the top 5 sales orders"
                  </div>
                </div>
              )}
              {messages.map((message) => (
                <div
                  key={message.id}
                  className={`flex ${message.role === "user" ? "justify-end" : "justify-start"}`}
                >
                  <div
                    className={`max-w-[80%] rounded-lg px-4 py-3 ${
                      message.role === "user"
                        ? "bg-primary text-primary-foreground"
                        : "bg-muted text-foreground"
                    }`}
                  >
                    <div className="whitespace-pre-wrap">{message.content}</div>
                    {/* Data table if present */}
                    {message.data && message.data.length > 0 && (
                      <div className="mt-3">
                        <details className="group">
                          <summary className="cursor-pointer text-sm font-medium underline underline-offset-4">
                            View data ({message.data.length} rows)
                          </summary>
                          <div className="mt-2 overflow-x-auto rounded border">
                            <table className="w-full text-sm">
                              <thead>
                                <tr className="border-b bg-muted/50">
                                  {Object.keys(message.data[0]).map((key) => (
                                    <th
                                      key={key}
                                      className="px-3 py-2 text-left font-medium"
                                    >
                                      {key}
                                    </th>
                                  ))}
                                </tr>
                              </thead>
                              <tbody>
                                {message.data.slice(0, 20).map((row, i) => (
                                  <tr key={i} className="border-b">
                                    {Object.entries(row).map(([key, value]) => (
                                      <td
                                        key={key}
                                        className="max-w-[200px] truncate px-3 py-2"
                                      >
                                        {value == null ? (
                                          <span className="text-muted-foreground">—</span>
                                        ) : typeof value === "object" ? (
                                          <span className="text-xs text-muted-foreground">
                                            {JSON.stringify(value)}
                                          </span>
                                        ) : (
                                          <span>{String(value)}</span>
                                        )}
                                      </td>
                                    ))}
                                  </tr>
                                ))}
                              </tbody>
                            </table>
                          </div>
                        </details>
                      </div>
                    )}
                    {/* Sources if present */}
                    {message.sources && message.sources.length > 0 && (
                      <div className="mt-2 text-xs text-muted-foreground">
                        Sources: {message.sources.join(", ")}
                      </div>
                    )}
                  </div>
                </div>
              ))}
              {mutation.isPending && (
                <div className="flex justify-start">
                  <div className="flex items-center gap-2 rounded-lg bg-muted px-4 py-3 text-muted-foreground">
                    <SpinnerGapIcon className="size-4 animate-spin" />
                    Thinking...
                  </div>
                </div>
              )}
            </div>
          </ScrollArea>

          {/* Input area */}
          <div className="border-t p-4">
            <form onSubmit={handleSubmit} className="flex gap-2">
              <Textarea
                ref={textareaRef}
                value={input}
                onChange={(e) => setInput(e.target.value)}
                onKeyDown={handleKeyDown}
                placeholder="Ask a question about your SAP data..."
                className="min-h-[44px] flex-1 resize-none"
                rows={1}
              />
              <Button
                type="submit"
                disabled={!input.trim() || mutation.isPending}
                className="h-[44px] w-[44px] p-0"
              >
                <PaperPlaneTiltIcon className="size-5" />
              </Button>
            </form>
          </div>
        </CardContent>
      </Card>
    </div>
  )
}

export const Route = createFileRoute("/ask/")({
  component: Ask,
})