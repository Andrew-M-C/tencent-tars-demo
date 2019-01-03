#include "CppWebServer.h"
#include "CppWebImp.h"

using namespace std;

CppWebServer g_app;

struct HttpProtocol
{
    /**
     * 解析http请求
     * @param in
     * @param out
     *
     * @return int
     */
    static int parseHttp(string &in, string &out)
    {
        try
        {
            //判断请求是否是HTTP请求
            bool b = TC_HttpRequest ::checkRequest(in.c_str(), in.length());
            //完整的HTTP请求
            if(b)
            {
                out = in;
                in  = "";
                //TLOGDEBUG("out size: " << out.size() << endl);
                return TC_EpollServer::PACKET_FULL;
            }
            else
            {
                return TC_EpollServer::PACKET_LESS;
            }
        }
        catch(exception &ex)
        {
            return TC_EpollServer::PACKET_ERR;
        }

        return TC_EpollServer::PACKET_LESS;             //表示收到的包不完全
    }

};

/////////////////////////////////////////////////////////////////
void
CppWebServer::initialize()
{
    //initialize application here:
    //...

    addServant<CppWebImp>(ServerConfig::Application + "." + ServerConfig::ServerName + ".CppWebObj");
	addServantProtocol(ServerConfig::Application + "." + ServerConfig::ServerName + ".CppWebObj", &HttpProtocol::parseHttp);
}
/////////////////////////////////////////////////////////////////
void
CppWebServer::destroyApp()
{
    //destroy application here:
    //...
}
/////////////////////////////////////////////////////////////////
int
main(int argc, char* argv[])
{
    try
    {
        g_app.main(argc, argv);
        g_app.waitForShutdown();
    }
    catch (std::exception& e)
    {
        cerr << "std::exception:" << e.what() << std::endl;
    }
    catch (...)
    {
        cerr << "unknown exception." << std::endl;
    }
    return -1;
}
/////////////////////////////////////////////////////////////////
