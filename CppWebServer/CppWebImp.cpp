#include "CppWebImp.h"
#include "servant/Application.h"

using namespace std;

//////////////////////////////////////////////////////
void CppWebImp::initialize()
{
    //initialize servant here:
    //...
}

//////////////////////////////////////////////////////
void CppWebImp::destroy()
{
    //destroy servant here:
    //...
}

int CppWebImp::doRequest(TarsCurrentPtr current, vector<char> &buffer)
{
    TC_HttpRequest req;
    TC_HttpResponse rsp;

	// parse request header
    vector<char> v = current->getRequestBuffer();
    string sBuf;
    sBuf.assign(&v[0], v.size());
    req.decode(sBuf);

    int ret = doRequest(req, rsp);

    rsp.encode(buffer);

    return ret;
}

int CppWebImp::doRequest(const TC_HttpRequest &req, TC_HttpResponse &rsp)
{
	string msg = "{\"msg\":\"Hello, Tars-Cpp!\"}";
    rsp.setContentType("application/json;charset=utf-8");
    rsp.setResponse(msg.c_str(), msg.size());
    return 0;
}
