Name:    ion-connect
Version: 0.7.3
Release: 1%{?dist}
Summary: CLI tool for accessing the Ion Channel API

License: ASL2.0
Source0: ion-connect
BuildArch: x86_64

%description
The official CLI tool for interacting with the Ion Channel API. This tool
allows Ion Channel users to access the full set of services provided by 
Ion Channel

%install
mkdir -p %{buildroot}/%{_bindir}
install -p -m 755 %{SOURCE0} %{buildroot}/%{_bindir}

%files
%{_bindir}/ion-connect

%changelog
