"use client";

import {
  Bold,
  Heading1,
  Heading2,
  Heading3,
  Italic,
  Link2,
  List,
  ListOrdered,
  Quote,
  Redo,
  Strikethrough,
  Undo,
} from "lucide-react";
import Image from "@tiptap/extension-image";
import Link from "@tiptap/extension-link";
import StarterKit from "@tiptap/starter-kit";
import { EditorContent, useEditor, type Editor } from "@tiptap/react";
import { useEffect, useRef } from "react";

import { Button } from "@/components/ui/button";
import { Separator } from "@/components/ui/separator";
import { cn } from "@/tools/utils";

interface TiptapEditorProps {
  value: string;
  onChange: (html: string) => void;
  placeholder?: string;
  disabled?: boolean;
}

export function TiptapEditor({
  value,
  onChange,
  disabled,
}: TiptapEditorProps) {
  // Track the HTML we've emitted so we don't bounce it back into setContent
  // and reset the user's cursor on every keystroke.
  const lastEmittedRef = useRef<string>(value || "");

  const editor = useEditor({
    extensions: [
      StarterKit.configure({
        heading: { levels: [1, 2, 3] },
      }),
      Link.configure({ openOnClick: false, autolink: true }),
      Image.configure({ inline: false }),
    ],
    content: value || "",
    editable: !disabled,
    immediatelyRender: false,
    // Tiptap v3 stopped re-rendering on every transaction by default. Without
    // this flag the toolbar's `editor.isActive(...)` / `editor.can().undo()`
    // checks never refresh, so active states + enabled states look frozen.
    shouldRerenderOnTransaction: true,
    editorProps: {
      attributes: {
        // Force a white canvas + dark text so the writing surface stays
        // readable regardless of the site theme. `dark:prose-invert` is
        // intentionally omitted for the same reason.
        class:
          "prose prose-sm prose-zinc max-w-none min-h-[280px] rounded-md border border-input bg-white text-zinc-900 caret-zinc-900 p-4 focus:outline-none focus-visible:ring-1 focus-visible:ring-ring",
      },
    },
    onUpdate: ({ editor }) => {
      const html = editor.getHTML();
      lastEmittedRef.current = html;
      onChange(html);
    },
  });

  // Only push external `value` changes into the editor (e.g. when switching
  // between articles, or programmatically resetting the form). If the change
  // came from us via onUpdate, lastEmittedRef will match and we skip the
  // setContent — preventing the cursor-jump-on-keystroke loop.
  useEffect(() => {
    if (!editor) return;
    const next = value || "";
    if (next !== lastEmittedRef.current) {
      lastEmittedRef.current = next;
      editor.commands.setContent(next, { emitUpdate: false });
    }
  }, [editor, value]);

  // Keep the editor's editable state in sync with the prop after mount —
  // useEditor() captures `editable` once during initialisation.
  useEffect(() => {
    editor?.setEditable(!disabled);
  }, [editor, disabled]);

  if (!editor) {
    return (
      <div className="min-h-[280px] rounded-md border border-input bg-background p-4 text-sm text-muted-foreground">
        Loading editor…
      </div>
    );
  }

  return (
    <div className="space-y-2">
      <Toolbar editor={editor} disabled={disabled} />
      <EditorContent editor={editor} />
    </div>
  );
}

interface ToolbarProps {
  editor: Editor;
  disabled?: boolean;
}

function Toolbar({ editor, disabled }: ToolbarProps) {
  return (
    <div className="flex flex-wrap items-center gap-1 rounded-md border bg-muted/40 p-1">
      <ToolbarButton
        label="Bold"
        active={editor.isActive("bold")}
        onClick={() => editor.chain().focus().toggleBold().run()}
        disabled={disabled}
      >
        <Bold className="h-4 w-4" />
      </ToolbarButton>
      <ToolbarButton
        label="Italic"
        active={editor.isActive("italic")}
        onClick={() => editor.chain().focus().toggleItalic().run()}
        disabled={disabled}
      >
        <Italic className="h-4 w-4" />
      </ToolbarButton>
      <ToolbarButton
        label="Strikethrough"
        active={editor.isActive("strike")}
        onClick={() => editor.chain().focus().toggleStrike().run()}
        disabled={disabled}
      >
        <Strikethrough className="h-4 w-4" />
      </ToolbarButton>
      <Separator orientation="vertical" className="mx-1 h-6" />
      <ToolbarButton
        label="Heading 1"
        active={editor.isActive("heading", { level: 1 })}
        onClick={() => editor.chain().focus().toggleHeading({ level: 1 }).run()}
        disabled={disabled}
      >
        <Heading1 className="h-4 w-4" />
      </ToolbarButton>
      <ToolbarButton
        label="Heading 2"
        active={editor.isActive("heading", { level: 2 })}
        onClick={() => editor.chain().focus().toggleHeading({ level: 2 }).run()}
        disabled={disabled}
      >
        <Heading2 className="h-4 w-4" />
      </ToolbarButton>
      <ToolbarButton
        label="Heading 3"
        active={editor.isActive("heading", { level: 3 })}
        onClick={() => editor.chain().focus().toggleHeading({ level: 3 }).run()}
        disabled={disabled}
      >
        <Heading3 className="h-4 w-4" />
      </ToolbarButton>
      <Separator orientation="vertical" className="mx-1 h-6" />
      <ToolbarButton
        label="Bullet list"
        active={editor.isActive("bulletList")}
        onClick={() => editor.chain().focus().toggleBulletList().run()}
        disabled={disabled}
      >
        <List className="h-4 w-4" />
      </ToolbarButton>
      <ToolbarButton
        label="Ordered list"
        active={editor.isActive("orderedList")}
        onClick={() => editor.chain().focus().toggleOrderedList().run()}
        disabled={disabled}
      >
        <ListOrdered className="h-4 w-4" />
      </ToolbarButton>
      <ToolbarButton
        label="Blockquote"
        active={editor.isActive("blockquote")}
        onClick={() => editor.chain().focus().toggleBlockquote().run()}
        disabled={disabled}
      >
        <Quote className="h-4 w-4" />
      </ToolbarButton>
      <Separator orientation="vertical" className="mx-1 h-6" />
      <ToolbarButton
        label="Link"
        active={editor.isActive("link")}
        onClick={() => {
          const previous = editor.getAttributes("link").href as string | undefined;
          const url = window.prompt("URL", previous ?? "https://");
          if (url === null) return;
          if (url === "") {
            editor.chain().focus().unsetLink().run();
            return;
          }
          editor.chain().focus().extendMarkRange("link").setLink({ href: url }).run();
        }}
        disabled={disabled}
      >
        <Link2 className="h-4 w-4" />
      </ToolbarButton>
      <Separator orientation="vertical" className="mx-1 h-6" />
      <ToolbarButton
        label="Undo"
        onClick={() => editor.chain().focus().undo().run()}
        disabled={disabled || !editor.can().undo()}
      >
        <Undo className="h-4 w-4" />
      </ToolbarButton>
      <ToolbarButton
        label="Redo"
        onClick={() => editor.chain().focus().redo().run()}
        disabled={disabled || !editor.can().redo()}
      >
        <Redo className="h-4 w-4" />
      </ToolbarButton>
    </div>
  );
}

function ToolbarButton({
  label,
  active,
  children,
  ...props
}: React.ButtonHTMLAttributes<HTMLButtonElement> & {
  label: string;
  active?: boolean;
}) {
  return (
    <Button
      type="button"
      variant={active ? "secondary" : "ghost"}
      size="icon"
      className={cn("h-8 w-8")}
      aria-label={label}
      title={label}
      {...props}
    >
      {children}
    </Button>
  );
}
