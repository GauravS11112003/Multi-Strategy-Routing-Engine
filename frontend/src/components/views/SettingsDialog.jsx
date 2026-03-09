import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogDescription } from "../ui/dialog";
import { Button } from "../ui/button";
import { Badge } from "../ui/badge";
import { ExternalLink, Key, Trash2, Check } from "lucide-react";

function SettingsDialog({ open, onOpenChange, apiKey, apiKeyInput, onApiKeyInputChange, onSave, onClear }) {
  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="max-w-md" onClose={() => onOpenChange(false)}>
        <DialogHeader>
          <DialogTitle className="flex items-center gap-2">
            <Key className="h-4 w-4 text-primary" />
            API Settings
          </DialogTitle>
          <DialogDescription>
            Configure your OpenRouteService API key for real road routing.
          </DialogDescription>
        </DialogHeader>

        <div className="space-y-4 mt-2">
          <div className="rounded-lg border border-border bg-secondary/30 p-4 space-y-3">
            <div className="flex items-center justify-between">
              <p className="text-xs font-medium">OpenRouteService API Key</p>
              {apiKey ? (
                <Badge variant="success">
                  <Check className="h-3 w-3 mr-1" />
                  Connected
                </Badge>
              ) : (
                <Badge variant="warning">Not Set</Badge>
              )}
            </div>

            <input
              type="password"
              value={apiKeyInput}
              onChange={(e) => onApiKeyInputChange(e.target.value)}
              placeholder="Paste your API key here..."
              className="w-full rounded-md border border-border bg-background px-3 py-2 text-sm font-mono placeholder:text-muted-foreground/50 focus:outline-none focus:ring-1 focus:ring-primary"
            />

            {apiKeyInput && (
              <p className="text-[10px] text-muted-foreground">
                {apiKeyInput.length} characters
              </p>
            )}

            <div className="flex gap-2">
              <Button
                size="sm"
                onClick={onSave}
                disabled={!apiKeyInput}
                className="flex-1"
              >
                Save Key
              </Button>
              {apiKey && (
                <Button size="sm" variant="outline" onClick={onClear}>
                  <Trash2 className="h-3.5 w-3.5" />
                </Button>
              )}
            </div>
          </div>

          <a
            href="https://openrouteservice.org/dev/#/signup"
            target="_blank"
            rel="noopener noreferrer"
            className="flex items-center gap-2 rounded-lg border border-border p-3 text-xs text-muted-foreground hover:text-foreground hover:bg-accent transition-colors"
          >
            <ExternalLink className="h-3.5 w-3.5" />
            <span>Get a free API key from OpenRouteService</span>
          </a>

          <div className="rounded-lg bg-secondary/30 p-3 border border-border">
            <p className="text-[10px] text-muted-foreground leading-relaxed">
              Your API key is stored locally in your browser and sent
              directly to OpenRouteService. It is never stored on any server.
            </p>
          </div>
        </div>
      </DialogContent>
    </Dialog>
  );
}

export default SettingsDialog;
