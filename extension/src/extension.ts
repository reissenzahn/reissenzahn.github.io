import {CancellationToken, CodeLens, ExtensionContext, languages, Position, Range, TextDocument, commands, window, ProgressLocation, Terminal} from 'vscode';
import {tmpdir} from 'os';
import {join} from 'path';
import {randomBytes} from 'crypto';
import {mkdirSync, writeFileSync} from 'fs';

const RUN_CODE_BLOCK_COMMAND = 'ext.runCodeBlock';
const CODE_BLOCK_DELIMITER = '```';
const TERMINAL_NAME = 'markdown-runner';

const provideCodeLenses = (document: TextDocument, cancellationToken: CancellationToken): CodeLens[] => {
  const codeLenses: CodeLens[] = [];

  for (let i = 0; i < document.lineCount; i++) {
    const line = document.lineAt(i);

    if (line.text.startsWith(CODE_BLOCK_DELIMITER) && line.text.substring(CODE_BLOCK_DELIMITER.length).trim()) {
      const type = line.text.substring(CODE_BLOCK_DELIMITER.length).trim();

      const start = i;

      let end;
      for (let j = i + 1; j < document.lineCount; j++) {
        if (document.lineAt(j).text.startsWith(CODE_BLOCK_DELIMITER)) {
          end = j;
          break;
        }
      }

      if (end && ['go'].includes(type)) {
        const content = document.getText(new Range(
          new Position(start + 1, 0),
          new Position(end - 1, document.lineAt(end - 1).range.end.character),
        ));

        codeLenses.push(new CodeLens(new Range(
          new Position(start, 0),
          new Position(start, 0),
        ), {
          command: RUN_CODE_BLOCK_COMMAND,
          title: 'Run',
          tooltip: `Run ${type} code block`,
          arguments: [{
            type,
            content,
          }]
        }))
      }
    }
  }

  return codeLenses;
}

let terminal: Terminal | null = null;

const runCodeBlock = async (type: string, content: string, context: ExtensionContext) => {
  await window.withProgress({
    location: ProgressLocation.Notification,
    title: 'Running code block',
    cancellable: true,
  }, async () => {
    const filePath = join(tmpdir(), 'markdown')
    mkdirSync(filePath, {recursive: true});

    const fileName = join(filePath, `${randomBytes(16).toString('hex')}.${type}`);
    writeFileSync(fileName, content, { mode: 0o600 });

    if (!terminal) {
      for (const t of window.terminals) {
        if (t.name == TERMINAL_NAME) {
          terminal = t;
          context.subscriptions.push(terminal);
          break;
        }
      }
    }

    if (!terminal) {
      terminal = window.createTerminal(TERMINAL_NAME);
      context.subscriptions.push(terminal);
    }

    terminal.show();
    commands.executeCommand('workbench.action.terminal.clear');

    if (type === 'go') {
      terminal.sendText(`go run "${fileName}"`, true);
    }
  })

  // TODO: cleanup temporary files
}

export const activate = (context: ExtensionContext) => {
  context.subscriptions.push(languages.registerCodeLensProvider({
    scheme: 'file',
    language: 'markdown',
  }, {
    provideCodeLenses,
  }));

  context.subscriptions.push(commands.registerCommand(RUN_CODE_BLOCK_COMMAND, (args) => {
    runCodeBlock(args.type, args.content, context);
  }))

  window.onDidCloseTerminal((t) => {
    if (t.name == TERMINAL_NAME) {
      terminal = null;
    }
  });
}

export const deactivate = () => {};
